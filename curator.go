package curator

import (
	"time"

	"github.com/casualjim/go-curator/ensemble"
	"github.com/obeattie/go-zookeeper/zk"
)

type RetryPolicy func(retryCount int, elapsed time.Duration) (shouldRetry bool)

func ForeverPolicy(retryCount int, elapsed time.Duration) (shouldRetry bool) {
	shouldRetry = true
	return
}

// Curator wraps a zookeeper connection and takes care of some house keeping
// to keep the connection reliable and so on.
type Curator struct {
	state             *ConnectionState
	RetryPolicy       RetryPolicy
	ConnectionTimeout time.Duration
}

// NewWithFactory creates a new curator with the provided factory and ensemble provider
func NewWithFactory(factory ZookeeperFactory, prov ensemble.EnsembleProvider, sessionTimeout time.Duration, connectionTimeout time.Duration) *Curator {
	return &Curator{
		state:             NewConnectionState(factory, prov, sessionTimeout, connectionTimeout),
		RetryPolicy:       ForeverPolicy,
		ConnectionTimeout: connectionTimeout,
	}
}

// NewFromUri creates a new curator for the connection string
func NewFromUri(uri string, sessionTimeout time.Duration, connectionTimeout time.Duration) (*Curator, error) {
	factory := DefaultFactory()
	prov, _, err := ensemble.New(uri)
	if err != nil {
		return nil, err
	}

	return NewWithFactory(factory, prov, sessionTimeout, connectionTimeout), nil
}

func (c *Curator) CurrentConnectionString() []string {
	return c.state.CurrentConnectionString()
}

// Close closes this zookeeper client, disconnects and cleans up state
func (c *Curator) Close() error {
	var err error
	if c.state != nil {
		err = c.state.Close()
	}
	return err
}

func (c *Curator) AddParentWatcher(watcher chan<- zk.Event) {
	c.state.AddParentWatcher(watcher)
}

func (c *Curator) RemoveParentWatcher(watcher chan<- zk.Event) {
	c.state.RemoveParentWatcher(watcher)
}

func (c *Curator) ConnectionIndex() int32 {
	return c.state.InstanceIndex()
}
