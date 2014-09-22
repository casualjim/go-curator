package curator

import (
	"errors"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/casualjim/go-curator/ensemble"
	"github.com/casualjim/go-curator/shared"
	"github.com/casualjim/go-zookeeper/zk"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("curator")

// RetryPolicy the predicate to determine if the operation should be retried after a failure
type RetryPolicy func(int, time.Duration) bool

// ForeverPolicy this policy keeps retrying forever without an exponential backoff.
func ForeverPolicy(retryCount int, elapsed time.Duration) bool {
	return true
}

type watcherFactory interface {
	MakeWatcher() chan zk.Event
	Invalidate()
}

type defaultWatcherFactory struct {
}

func (d *defaultWatcherFactory) MakeWatcher() chan zk.Event {
	return make(chan zk.Event)
}

func (d *defaultWatcherFactory) Invalidate() {

}

// Conn wraps a zookeeper connection and takes care of some house keeping
// to keep the connection reliable and so on.
type Conn struct {
	sync.Mutex
	state             ConnectionState
	RetryPolicy       RetryPolicy
	ConnectionTimeout time.Duration
	started           int32
	startedPtr        *int32
	watcherFactory    watcherFactory
}

// NewWithFactory creates a new Conn with the provided factory and ensemble provider
func NewWithFactory(factory ZookeeperFactory, prov ensemble.Provider, sessionTimeout time.Duration, connectionTimeout time.Duration) (*Conn, error) {
	if sessionTimeout < connectionTimeout {
		logger.Warning("session timeout [%d] is less than connection timeout [%d]", sessionTimeout, connectionTimeout)
	}

	conn := &Conn{
		state:             newConnectionState(factory, prov, sessionTimeout, connectionTimeout),
		RetryPolicy:       ForeverPolicy,
		ConnectionTimeout: connectionTimeout,
		started:           0,
		watcherFactory:    &defaultWatcherFactory{},
	}
	conn.startedPtr = &conn.started
	return conn, nil
}

// NewFromURI creates a new Conn for the connection string
func NewFromURI(uri string, sessionTimeout time.Duration, connectionTimeout time.Duration) (*Conn, error) {
	factory := DefaultZookeeperFactory()
	prov, _, err := ensemble.Fixed(uri)
	if err != nil {
		return nil, err
	}

	return NewWithFactory(factory, prov, sessionTimeout, connectionTimeout)
}

// CurrentConnectionString the hosts for the current connection
func (c *Conn) CurrentConnectionString() (string, error) {
	if c.started == 0 {
		return "", errors.New("The client needs to be started before you can get a connection string")
	}
	return strings.Join(c.state.CurrentConnectionString(), ","), nil
}

// Start starts this curator connection
func (c *Conn) Start() error {
	c.Lock()
	defer c.Unlock()
	if c.started == 1 {
		logger.Warning("Called start on Conn more than once")
		return nil
	}
	err := c.state.Start()
	if err != nil {
		return err
	}
	c.started = 1
	return nil
}

// Close closes this zookeeper client, disconnects and cleans up state
func (c *Conn) Close() error {
	c.Lock()
	defer c.Unlock()
	c.started = 0
	var err error
	if c.state != nil {
		err = c.state.Close()
	}
	return err
}

// AddParentWatcher adds a connection watcher, receives zookeeper events
func (c *Conn) AddParentWatcher(watcher chan<- zk.Event) {
	c.state.AddParentWatcher(watcher)
}

// RemoveParentWatcher removes a connection watcher
func (c *Conn) RemoveParentWatcher(watcher chan<- zk.Event) {
	c.state.RemoveParentWatcher(watcher)
}

// ConnectionIndex the index of this connection, the amount of reconnections
func (c *Conn) ConnectionIndex() int32 {
	return c.state.InstanceIndex()
}

// IsConnected() returns true when this client is connected
func (c *Conn) IsConnected() bool {
	logger.Debug("calling is connected")
	return c.state.IsConnected()
}

// BlockUntilConnectedOrTimedOut blocks until the connection to ZK succeeds. Use with caution. The block
// will timeout after the connection timeout (as passed to the constructor) has elapsed
func (c *Conn) BlockUntilConnectedOrTimedOut() (bool, error) {
	if c.started == 0 {
		return false, errors.New("The client needs to be started before you can make a connection")
	}
	logger.Debug("BlockUntilConnectedOrTimedOut start")

	c.internalBlockUntilConnectedOrTimedOut()

	isConnected := c.IsConnected()
	logger.Debug("BlockUntilConnectedOrTimedOut end. isConnected: %t", isConnected)
	return isConnected, nil

}

func (c *Conn) internalBlockUntilConnectedOrTimedOut() {
	logger.Debug("Entering internalBlockUntilConnectedOrTimedOut")
	waitTime := c.ConnectionTimeout.Nanoseconds()
	logger.Debug("about to start loop")
	for !c.IsConnected() && waitTime > 0 {
		logger.Debug("Passing through internalBlockUntilConnectedOrTimedOut loop")
		// waitTime = c.withTempWatcher(waitTime, func(watcher chan zk.Event) {
		// 	select { // Block until timeout or until a connection event was received
		// 	case <-watcher:
		// 	case <-time.After(1 * time.Second):
		// 	}
		// })
		// watcher := c.watcherFactory.MakeWatcher()
		watcher := &shared.WatcherHolder{make(chan zk.Event)}
		c.state.AddParentWatcherHolder(watcher)
		startTime := time.Now()

		select { // Block until timeout or until a connection event was received
		case <-watcher.Watcher:
			logger.Debug("the watcher received an event")
		case <-time.After(1 * time.Second):
			logger.Debug("this loop timed out")
		}
		logger.Debug("Passed the select block")
		c.state.RemoveParentWatcherHolder(watcher)
		// c.watcherFactory.Invalidate()
		elapsed := math.Max(1, float64(time.Now().UnixNano()-startTime.UnixNano()))
		waitTime = waitTime - int64(elapsed)
		logger.Debug("exiting loop")
	}
	logger.Debug("Leaving internalBlockUntilConnectedOrTimedOut")
}

func (c *Conn) withTempWatcher(waitTime int64, thunk func(watcher chan zk.Event)) int64 {
	watcher := make(chan zk.Event)

	c.AddParentWatcher(watcher)
	startTime := time.Now()

	thunk(watcher)

	c.RemoveParentWatcher(watcher)
	elapsed := math.Max(1, float64(time.Now().UnixNano()-startTime.UnixNano()))
	return waitTime - int64(elapsed)
}
