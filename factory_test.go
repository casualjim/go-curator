package curator

import (
	"fmt"
	"testing"
	"time"

	"github.com/casualjim/go-zookeeper/zk"
	. "github.com/onsi/gomega"
	. "github.com/smartystreets/goconvey/convey"
)

func TestZookeeperFactory(t *testing.T) {
	Convey("A ZookeeperFactory should connect to a zookeeper ensemble", t, func() {

		Convey("connect to a zookeeper ensemble", func() {
			zkCluster, err := zk.StartTestCluster(1)
			So(err, ShouldBeNil)
			defer zkCluster.Stop()
			serv := zkCluster.Servers[0]
			hosts := []string{fmt.Sprintf("localhost:%d", serv.Port)}
			factory := DefaultZookeeperFactory()
			c, _, err := factory.NewZookeeper(hosts, 1*time.Second)
			defer c.Close()
			So(err, ShouldBeNil)
			Eventually(func() zk.State {
				return c.State()
			}).Should(Equal(zk.StateHasSession))
			So(c.State(), ShouldEqual, zk.StateHasSession)
		})

		Convey("eventually receive a Session event and have a sesssion", func() {
			zkCluster, err := zk.StartTestCluster(1)
			So(err, ShouldBeNil)
			defer zkCluster.Stop()
			serv := zkCluster.Servers[0]
			hosts := []string{fmt.Sprintf("localhost:%d", serv.Port)}
			factory := DefaultZookeeperFactory()
			c, watcher, err := factory.NewZookeeper(hosts, 1*time.Second)
			defer c.Close()
			So(err, ShouldBeNil)

			var event zk.Event
			Eventually(func() zk.State {
				event = <-watcher
				fmt.Println("The current event:", event)
				return event.State
			}).Should(Equal(zk.StateHasSession))
			So(event.State, ShouldEqual, zk.StateHasSession)
			So(event.Type, ShouldEqual, zk.EventSession)
			So(event.Path, ShouldEqual, "")
			So(event.Err, ShouldBeNil)
		})

	})
}
