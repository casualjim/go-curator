package curator

import (
	"time"

	"github.com/obeattie/go-zookeeper/zk"
)

type ZookeeperFactory interface {
	NewZookeeper(hosts []string, timeout time.Duration) (zk.IConn, <-chan zk.Event, error)
}

type defaultZookeeperFactory struct {
}

func (f *defaultZookeeperFactory) NewZookeeper(servers []string, timeout time.Duration) (zk.IConn, <-chan zk.Event, error) {
	return zk.Connect(servers, timeout)
}

func DefaultFactory() ZookeeperFactory {
	return &defaultZookeeperFactory{}
}
