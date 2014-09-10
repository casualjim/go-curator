package curator

import (
	"time"

	"github.com/casualjim/go-zookeeper/zk"
)

// ZookeeperFactory an abstraction over connecting to zookeeper, useful in tests
type ZookeeperFactory interface {
	NewZookeeper(hosts []string, sessionTimeout time.Duration) (zk.IConn, <-chan zk.Event, error)
}

type defaultZookeeperFactory struct {
}

func (f *defaultZookeeperFactory) NewZookeeper(servers []string, sessionTimeout time.Duration) (zk.IConn, <-chan zk.Event, error) {
	logger.Debug("Connecting to %v", servers)
	return zk.Connect(servers, sessionTimeout)
}

// DefaultZookeeperFactory returns the default zookeeper factory that makes an actual connection
// to the specified ensemble hosts.
func DefaultZookeeperFactory() ZookeeperFactory {
	return &defaultZookeeperFactory{}
}
