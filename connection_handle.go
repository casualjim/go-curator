package curator

import (
	"log"
	"time"

	"github.com/casualjim/go-curator/ensemble"
	"github.com/obeattie/go-zookeeper/zk"
)

type connHandle struct {
	factory          ZookeeperFactory
	ensembleProvider ensemble.EnsembleProvider
	sessionTimeout   time.Duration
	current          *currentConnection
}

func (c *connHandle) HasNewConnectionString() bool {
	if c.current != nil {
		hc, l := c.current.hosts, len(c.current.hosts)
		prov := c.ensembleProvider.Hosts()
		if l > 0 && l == len(prov) { // compare host collections
			for _, h := range hc {
				var found bool
				for _, p := range prov {
					if h == p {
						found = true
						break
					}
				}
				if !found {
					return true
				}
			}
		} else { // only when the if guard doesn't match
			return l > 0
		}
	}
	return false
}

func (c *connHandle) Conn() zk.IConn {
	if c.current != nil {
		return c.current.conn
	}
	return nil
}

func (c *connHandle) Close() (err error) {
	err = c.internalClose()
	c.current = nil
	return
}

func (c *connHandle) Reconnect() (<-chan zk.Event, error) {
	if err := c.internalClose(); err != nil {
		return nil, err
	}
	hosts := c.ensembleProvider.Hosts()
	conn, w, err := c.factory.NewZookeeper(hosts, c.sessionTimeout)
	if err != nil {
		return nil, err
	}
	c.current = &currentConnection{hosts: hosts, conn: conn}
	return w, nil
}

func (c *connHandle) internalClose() error {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()
	if conn := c.Conn(); conn != nil {
		conn.Close()
	}
	return nil
}

type currentConnection struct {
	hosts []string
	conn  zk.IConn
}
