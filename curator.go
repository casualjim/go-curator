package curator

import (
	"time"

	"github.com/casualjim/go-curator/ensemble"
	"github.com/casualjim/go-zookeeper/zk"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("curator")

// RetryPolicy the predicate to determine if the operation should be retried after a failure
type RetryPolicy func(retryCount int, elapsed time.Duration) (shouldRetry bool)

// ForeverPolicy this policy keeps retrying forever without an exponential backoff.
func ForeverPolicy(retryCount int, elapsed time.Duration) (shouldRetry bool) {
	shouldRetry = true
	return
}

// CuratorConn wraps a zookeeper connection and takes care of some house keeping
// to keep the connection reliable and so on.
type CuratorConn struct {
	state             *connectionState
	RetryPolicy       RetryPolicy
	ConnectionTimeout time.Duration
}

// NewWithFactory creates a new CuratorConn with the provided factory and ensemble provider
func NewWithFactory(factory ZookeeperFactory, prov ensemble.Provider, sessionTimeout time.Duration, connectionTimeout time.Duration) *CuratorConn {
	return &CuratorConn{
		state:             newConnectionState(factory, prov, sessionTimeout, connectionTimeout),
		RetryPolicy:       ForeverPolicy,
		ConnectionTimeout: connectionTimeout,
	}
}

// NewFromURI creates a new CuratorConn for the connection string
func NewFromURI(uri string, sessionTimeout time.Duration, connectionTimeout time.Duration) (*CuratorConn, error) {
	factory := DefaultZookeeperFactory()
	prov, _, err := ensemble.Fixed(uri)
	if err != nil {
		return nil, err
	}

	return NewWithFactory(factory, prov, sessionTimeout, connectionTimeout), nil
}

// CurrentConnectionString the hosts for the current connection
func (c *CuratorConn) CurrentConnectionString() []string {
	return c.state.CurrentConnectionString()
}

// Close closes this zookeeper client, disconnects and cleans up state
func (c *CuratorConn) Close() error {
	var err error
	if c.state != nil {
		err = c.state.Close()
	}
	return err
}

// AddParentWatcher adds a connection watcher, receives zookeeper events
func (c *CuratorConn) AddParentWatcher(watcher chan<- zk.Event) {
	c.state.AddParentWatcher(watcher)
}

// RemoveParentWatcher removes a connection watcher
func (c *CuratorConn) RemoveParentWatcher(watcher chan<- zk.Event) {
	c.state.RemoveParentWatcher(watcher)
}

// ConnectionIndex the index of this connection, roughly the amount of reconnections
func (c *CuratorConn) ConnectionIndex() int32 {
	return c.state.InstanceIndex()
}
