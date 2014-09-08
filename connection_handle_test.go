// Package curator provides ...
package curator

import (
	"testing"
	"time"

	"code.google.com/p/gomock/gomock"

	"github.com/casualjim/go-curator/ensemble"
	"github.com/casualjim/go-curator/mock_zookeeper"
	"github.com/casualjim/go-zookeeper/zk"
	. "github.com/smartystreets/goconvey/convey"
)

type mockZookeeperFactory struct {
	mockCtrl  *gomock.Controller
	channel   chan zk.Event
	conn      *mock_zookeeper.MockIConn
	callCount int
}

func (f *mockZookeeperFactory) NewZookeeper(hosts []string, timeout time.Duration) (zk.IConn, <-chan zk.Event, error) {
	f.callCount++
	f.conn = mock_zookeeper.NewMockIConn(f.mockCtrl)
	return f.conn, f.channel, nil
}

func testFactory(mocks *gomock.Controller) *mockZookeeperFactory {
	return &mockZookeeperFactory{mockCtrl: mocks, channel: make(chan zk.Event)}
}

func TestConnHandle(t *testing.T) {
	Convey("A connection handle should", t, func() {
		mockCtrl := gomock.NewController(t)
		Reset(func() {
			mockCtrl.Finish()
		})

		Convey("return nil for the conn when there is no active connection", func() {
			fact := testFactory(mockCtrl)
			prov, _, _ := ensemble.Fixed("zk://localhost:2181/blah")
			handle := &connHandle{factory: fact, ensembleProvider: prov, sessionTimeout: 1 * time.Second}
			So(handle.Conn(), ShouldBeNil)
		})
		Convey("when disconnected", func() {

			fact := testFactory(mockCtrl)
			prov, _, _ := ensemble.Fixed("zk://localhost:2181/blah")
			handle := &connHandle{factory: fact, ensembleProvider: prov, sessionTimeout: 1 * time.Second}
			Convey("should connect", func() {
				_, err := handle.Reconnect()

				So(err, ShouldBeNil)
				So(handle.Conn(), ShouldNotBeNil)
			})

			Convey("not have a new connection string when disconnected", func() {
				So(handle.HasNewConnectionString(), ShouldBeFalse)
			})
		})

		Convey("when connected", func() {
			fact := testFactory(mockCtrl)
			prov, _, _ := ensemble.Fixed("zk://localhost:2181/blah")
			handle := &connHandle{factory: fact, ensembleProvider: prov, sessionTimeout: 1 * time.Second}
			_, err := handle.Reconnect()
			So(err, ShouldBeNil)
			conn := handle.Conn
			So(conn, ShouldNotBeNil)

			Convey("not have a new connection string when the list of hosts hasn't changed", func() {
				So(handle.HasNewConnectionString(), ShouldBeFalse)
			})

			Convey("have a new connection string when the list hosts is different in length", func() {
				prov, _, _ = ensemble.Fixed("zk://localhost:2181,remote-blah:2181/blah")
				handle.ensembleProvider = prov
				So(handle.HasNewConnectionString(), ShouldBeTrue)
			})

			Convey("have a new connection string when the host is different", func() {
				prov, _, _ = ensemble.Fixed("zk://remote-blah:2181/blah")
				handle.ensembleProvider = prov
				So(handle.HasNewConnectionString(), ShouldBeTrue)
			})

			Convey("close the connection before reconnecting", func() {
				fact.conn.EXPECT().Close()
				handle.Reconnect()
			})
		})
	})
}
