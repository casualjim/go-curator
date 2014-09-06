// Package curator provides reliable zookeeper client that keeps track of changes in the ensemble
package curator

import (
	"container/list"
	"fmt"
	"math"
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
	isConnected       int32
	conn              connHandle
	parentWatchers    events.EventBus
	started           time.Time
	connectionCount   int32
	connectionTimeout time.Duration
	errorQueue        *list.List
}

// NewConnectionState creates a new instance of connection state for the provided params
func NewConnectionState(factory ZookeeperFactory, ensembleProvider ensemble.EnsembleProvider, sessionTimeout time.Duration, connectionTimeout time.Duration) *ConnectionState {
	return &ConnectionState{
		parentWatchers:    events.New(),
		conn:              connHandle{factory: factory, ensembleProvider: ensembleProvider, sessionTimeout: sessionTimeout},
		connectionTimeout: connectionTimeout,
		errorQueue:        list.New(),
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

func (c *ConnectionState) InstanceIndex() int32 {
	return c.connectionCount
}

func (c *ConnectionState) reset() error {
	c.connectionCount++
	c.started = time.Now()
	c.isConnected = 0
	w, err := c.conn.Reconnect()
	if err != nil {
		return err
	}
	for {
		e := <-w
		if e.State == zk.StateHasSession {
			break
		} else if e.State == zk.StateExpired {
			return fmt.Errorf("Connection to %q expired %d.", c.conn.ensembleProvider.Hosts(), c.conn.sessionTimeout)
		}
	}
	c.connectionLoop(w)
	return nil
}

func (c *ConnectionState) CurrentConnectionString() []string {
	return c.conn.ensembleProvider.Hosts()
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
	err := c.reset()
	if err != nil {
		if c.errorQueue.Len() >= 10 {
			el := c.errorQueue.Back()
			c.errorQueue.Remove(el)
		}
		c.errorQueue.PushFront(err)
		return
	}
}

func (c *ConnectionState) Conn() (zk.IConn, error) {
	var err error = nil
	for e := c.errorQueue.Front(); e != nil; e = e.Next() {
		err = e.Value.(error)
		c.errorQueue.Remove(e)
	}
	if err != nil {
		return nil, err
	}
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

func (c *ConnectionState) checkTimeouts() error {
	sTo, cTo := float64(c.conn.sessionTimeout.Nanoseconds()), float64(c.connectionTimeout.Nanoseconds())
	minTimeout := int64(math.Min(sTo, cTo))
	elapsed := time.Now().UnixNano() - minTimeout

	if elapsed >= minTimeout {
		if c.conn.HasNewConnectionString() {
			c.backgroundReset()
		} else {

			maxTimeout := int64(math.Max(sTo, cTo))

			if elapsed > maxTimeout {
				return c.reset()
			} else {
				return fmt.Errorf("Curator connection timed out with %v/%v", maxTimeout, elapsed)
			}
		}
	}
	return nil
}
