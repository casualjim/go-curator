package curator

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"code.google.com/p/gomock/gomock"

	"github.com/casualjim/go-curator/ensemble"
	"github.com/casualjim/go-curator/mocks"
	"github.com/casualjim/go-zookeeper/zk"
	. "github.com/smartystreets/goconvey/convey"
)

func ensembleFromCluster(zkCluster *zk.TestCluster) ensemble.Provider {
	hosts := make([]string, len(zkCluster.Servers))
	for _, server := range zkCluster.Servers {
		hosts = append(hosts, fmt.Sprintf("localhost:%d", server.Port))
	}
	prov, _, _ := ensemble.Fixed(strings.Join(hosts, ","))
	return prov
}

func TestCurator(t *testing.T) {

	Convey("The curator zookeeper client should", t, func() {
		//RegisterTestingT(t)
		mockCtrl := gomock.NewController(t)

		factory := mocks.NewMockZookeeperFactory(mockCtrl)
		provider := mocks.NewMockEnsembleProvider(mockCtrl)

		conn, _ := NewWithFactory(factory, provider, 10*time.Second, 3*time.Second)

		Reset(func() {
			conn.Close()
			mockCtrl.Finish()
		})

		//zkCluster, err := zk.StartTestCluster(1)
		//So(err, ShouldBeNil)
		//factory := DefaultZookeeperFactory()
		//provider, _, _ := ensemble.Fixed(fmt.Sprintf("localhost:%d", zkCluster.Servers[0].Port))
		//conn, _ := NewWithFactory(factory, provider, 10*time.Second, 3*time.Second)
		//conn.Start()

		//Reset(func() {
		//conn.Close()
		//zkCluster.Stop()
		//})

		//Convey("connect to the cluster", func() {
		//success, err := conn.BlockUntilConnectedOrTimedOut()
		//So(err, ShouldBeNil)
		//So(success, ShouldBeTrue)
		//})
	})

}
