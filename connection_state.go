// Package curator provides reliable zookeeper client that keeps track of changes in the ensemble
package curator

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/casualjim/go-curator/ensemble"
	"github.com/casualjim/go-curator/shared/events"
	"github.com/obeattie/go-zookeeper/zk"
)

// ConnectionState implements the reliable connection to zookeeper
// and raising the events when the state in the connection changes.
// Should it detect that the connection string changes, it will disconnect the previous
// connection and reconnect with the newly detected ensemble
type ConnectionState struct {
	sync.Mutex
	isConnected     int32
	conn            connHandle
	parentWatchers  events.EventBus
	started         time.Time
	connectionCount int32
}

// NewConnectionState creates a new instance of connection state for the provided params
func NewConnectionState(factory ZookeeperFactory, ensembleProvider ensemble.EnsembleProvider, sessionTimeout time.Duration, connectionTimeout time.Duration) *ConnectionState {
	return &ConnectionState{
		parentWatchers: events.New(),
		conn:           connHandle{factory: factory, ensembleProvider: ensembleProvider, sessionTimeout: sessionTimeout},
	}
}

// AddParentWatcher adds a connection watcher, gets notified when something happens with the connection
func (c *ConnectionState) AddParentWatcher(watcher chan<- zk.Event) {
	c.parentWatchers.Add(watcher)
}

// RemoveParentWatcher removes a parent watcher
func (c *ConnectionState) RemoveParentWatcher(watcher chan<- zk.Event) {
	c.parentWatchers.Remove(watcher)
}

// IsConnected returns true when this connection is actually connected
func (c *ConnectionState) IsConnected() bool {
	return c.isConnected > 0
}

func (c *ConnectionState) Start() error {
	c.Lock()
	defer c.Unlock()
	return c.reset()
}

func (c *ConnectionState) Close() error {
	if atomic.CompareAndSwapInt32(&c.isConnected, 1, 0) {
		return c.conn.Close()

	}
	return nil
}

func (c *ConnectionState) reset() error {
	c.connectionCount++
	c.started = time.Now()
	c.isConnected = 0
	w, err := c.conn.Reconnect()
	if err != nil {
		return err
	}
	c.connectionLoop(w)
	return nil
}

func (c *ConnectionState) connectionLoop(connWatch <-chan zk.Event) {
	go func() {
		for evt := range connWatch {
			c.parentWatchers.Trigger() <- evt
			wasConnected := c.IsConnected()
			newIsConnected := wasConnected
			if evt.Type == zk.EventSession {
				newIsConnected = c.checkEvent(evt, wasConnected)
			}
			if newIsConnected != wasConnected {
				c.Lock()
				c.isConnected = 1
				c.started = time.Now()
				c.Unlock()
			}
		}
	}()
}

func (c *ConnectionState) backgroundReset() {
	c.Lock()
	defer c.Unlock()
	// TODO: Put this error on a queue, and use in the Conn() method
	c.reset()
}

func (c *ConnectionState) Conn() (zk.IConn, error) {
	return c.conn.Conn(), nil
}

func (c *ConnectionState) checkEvent(evt zk.Event, wasConnected bool) bool {
	isConnected := wasConnected
	checkNewConnectionString := true
	switch evt.State {
	case zk.StateDisconnected:
		isConnected = false
	case zk.StateSyncConnected | zk.StateConnectedReadOnly:
		isConnected = true
	case zk.StateAuthFailed:
		isConnected = false
	case zk.StateExpired:
		isConnected = false
		checkNewConnectionString = false
		c.backgroundReset()
	case zk.StateSaslAuthenticated:
	default:
		isConnected = false
	}
	if checkNewConnectionString && c.conn.HasNewConnectionString() {
		c.backgroundReset()
	}
	return isConnected
}
