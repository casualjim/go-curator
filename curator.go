package curator

import (
	"time"

	"github.com/casualjim/go-curator/ensemble"
	"github.com/obeattie/go-zookeeper/zk"
)

type RetryPolicy func(int, time.Duration) bool

type connectionState struct {
	factory           ZookeeperFactory
	ensembleProvider  ensemble.EnsembleProvider
	sessionTimeout    time.Duration
	connectionTimeout time.Duration
	connWatch         chan<- zk.Event
}

func (c *connectionState) Close() error {
	return nil
}

// Curator wraps a zookeeper connection and takes care of some house keeping
// to keep the connection reliable and so on.
type Curator struct {
	state *connectionState
}

func NewWithFactory(factory ZookeeperFactory, sessionTimeout time.Duration, connectionTimeout time.Duration) *Curator {
	return &Curator{state: &connectionState{factory: factory}}
}

func NewFromUri(uri string, sessionTimeout time.Duration, connectionTimeout time.Duration) *Curator {
	return nil
}

// Close closes this zookeeper client, disconnects and cleans up state
func (c *Curator) Close() error {
	var err error
	if c.state != nil {
		err = c.state.Close()
	}
}
