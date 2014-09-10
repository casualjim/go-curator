// Package curator provides reliable zookeeper client that keeps track of changes in the ensemble
package curator

import (
	"container/list"
	"fmt"
	"io"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/casualjim/go-curator/ensemble"
	"github.com/casualjim/go-curator/shared/events"
	"github.com/casualjim/go-zookeeper/zk"
	"github.com/rcrowley/go-metrics"
)

type ParentWatchable interface {
	AddParentWatcher(watcher chan<- zk.Event)
	RemoveParentWatcher(watcher chan<- zk.Event)
}

type Startable interface {
	Start() error
}

type ConnectionState interface {
	io.Closer
	ParentWatchable
	Startable

	InstanceIndex() int32
	CurrentConnectionString() []string
	Conn() (zk.IConn, error)
	IsConnected() bool
}

// connState implements the reliable connection to zookeeper
// and raising the events when the state in the connection changes.
// Should it detect that the connection string changes, it will disconnect the previous
// connection and reconnect with the newly detected ensemble
type connState struct {
	sync.Mutex
	isConnected       int32
	conn              ConnectionHandle
	parentWatchers    events.EventBus
	started           time.Time
	connectionCount   int32
	connectionTimeout time.Duration
	errorQueue        *list.List
}

// NewConnectionState creates a new instance of connection state for the provided params
func newConnectionState(factory ZookeeperFactory, ensembleProvider ensemble.Provider, sessionTimeout time.Duration, connectionTimeout time.Duration) *connState {
	return &connState{
		parentWatchers:    events.New(),
		conn:              &connHandle{factory: factory, ensembleProvider: ensembleProvider, sessionTimeout: sessionTimeout},
		connectionTimeout: connectionTimeout,
		errorQueue:        list.New(),
	}
}

// AddParentWatcher adds a connection watcher, gets notified when something happens with the connection
func (c *connState) AddParentWatcher(watcher chan<- zk.Event) {
	c.parentWatchers.Add(watcher)
}

// RemoveParentWatcher removes a parent watcher
func (c *connState) RemoveParentWatcher(watcher chan<- zk.Event) {
	c.parentWatchers.Remove(watcher)
}

// IsConnected returns true when this connection is actually connected
func (c *connState) IsConnected() bool {
	return c.isConnected > 0
}

func (c *connState) Start() error {
	c.Lock()
	defer c.Unlock()
	return c.reset()
}

func (c *connState) Close() error {
	if atomic.CompareAndSwapInt32(&c.isConnected, 1, 0) {
		return c.conn.Close()

	}
	return nil
}

func (c *connState) InstanceIndex() int32 {
	return c.connectionCount
}

func (c *connState) reset() error {
	c.connectionCount++
	c.started = time.Now()
	wasConnected := c.IsConnected()
	c.isConnected = 0
	w, err := c.conn.Reconnect()
	if err != nil {
		return err
	}
	c.connectionLoop(w, wasConnected)
	return nil
}

func (c *connState) CurrentConnectionString() []string {
	return c.conn.Hosts()
}

func (c *connState) connectionLoop(connWatch <-chan zk.Event, wasConnected bool) {
	go func() {
		for evt := range connWatch {
			c.parentWatchers.Trigger() <- evt
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

func (c *connState) handleNewConnectionString() {
	logger.Info("Connection string changed")
	metrics.GetOrRegisterCounter("curator.connection.hosts.changed", metrics.DefaultRegistry).Inc(1)
	c.backgroundReset()
}

func (c *connState) handleSessionExpired() {
	logger.Warning("Session expired event received")
	metrics.GetOrRegisterCounter("curator.connection.session.expired", metrics.DefaultRegistry).Inc(1)
	c.backgroundReset()
}

func (c *connState) backgroundReset() {
	c.Lock()
	defer c.Unlock()
	err := c.reset()
	if err != nil {
		if c.errorQueue.Len() >= 10 {
			el := c.errorQueue.Front()
			c.errorQueue.Remove(el)
		}
		c.errorQueue.PushBack(err)
		return
	}
}

func (c *connState) Conn() (zk.IConn, error) {
	var err error
	for e := c.errorQueue.Front(); e != nil; e = e.Next() {
		metrics.GetOrRegisterCounter("curator.connection.background-error", metrics.DefaultRegistry).Inc(1)
		err = e.Value.(error)
		c.errorQueue.Remove(e)
	}
	if err != nil {
		return nil, err
	}
	if c.IsConnected() {
		err = c.checkTimeouts()
		if err != nil {
			return nil, err
		}
	}
	return c.conn.Conn(), nil
}

func (c *connState) checkEvent(evt zk.Event, wasConnected bool) bool {
	isConnected := wasConnected
	checkNewConnectionString := true
	switch evt.State {
	case zk.StateDisconnected:
		isConnected = false
	case zk.StateSyncConnected:
		fallthrough
	case zk.StateConnectedReadOnly:
		isConnected = true
	case zk.StateAuthFailed:
		logger.Critical("Authentication failed")
		isConnected = false
	case zk.StateExpired:
		isConnected = false
		checkNewConnectionString = false
		c.handleSessionExpired()
	case zk.StateSaslAuthenticated:
	default:
		isConnected = false
	}

	if checkNewConnectionString && c.conn.HasNewConnectionString() {
		c.handleNewConnectionString()
	}
	return isConnected
}

func (c *connState) checkTimeouts() error {
	sTo, cTo := float64(c.conn.SessionTimeout().Nanoseconds()), float64(c.connectionTimeout.Nanoseconds())
	minTimeout := int64(math.Min(sTo, cTo))
	elapsed := time.Now().UnixNano() - c.started.UnixNano()

	if elapsed >= minTimeout {
		if c.conn.HasNewConnectionString() {
			c.handleNewConnectionString()
		} else {
			maxTimeout := int64(math.Max(sTo, cTo))
			if elapsed > maxTimeout {
				return c.reset()
			} else {
				metrics.GetOrRegisterCounter("curator.connection.timeout", metrics.DefaultRegistry).Inc(1)
				return fmt.Errorf("Curator connection timed out with %v/%v", maxTimeout, elapsed)
			}
		}
	}
	return nil
}
