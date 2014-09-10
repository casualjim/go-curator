package curator

import (
	"testing"
	"time"

	"code.google.com/p/gomock/gomock"

	"github.com/casualjim/go-curator/mocks"
	"github.com/casualjim/go-zookeeper/zk"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCurator(t *testing.T) {

	Convey("The curator zookeeper client should", t, func() {
		//RegisterTestingT(t)
		mockCtrl := gomock.NewController(t)
		state := mocks.NewMockConnectionState(mockCtrl)
		conn := &CuratorConn{
			state:             state,
			RetryPolicy:       ForeverPolicy,
			ConnectionTimeout: 3 * time.Second,
			started:           0,
		}
		conn.startedPtr = &conn.started

		Reset(func() {
			conn.Close()
			mockCtrl.Finish()
		})

		Convey("when not yet started", func() {
			state.EXPECT().Close().AnyTimes()

			Convey("should set started to 1", func() {
				state.EXPECT().Start().Times(1).Return(nil)
				res := conn.Start()
				So(res, ShouldBeNil)

				So(conn.started, ShouldEqual, 1)
			})

			Convey("getting the connection string returns an error", func() {
				str, err := conn.CurrentConnectionString()
				So(str, ShouldBeBlank)
				So(err, ShouldNotBeNil)
			})

			Convey("BlockUntilConnectedOrTimedOut returns an error", func() {
				connected, err := conn.BlockUntilConnectedOrTimedOut()
				So(connected, ShouldBeFalse)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("when already started", func() {
			state.EXPECT().Start().Times(1)
			state.EXPECT().Close().AnyTimes()
			conn.Start()

			Convey("should not try to start a second connection", func() {
				res := conn.Start()
				So(res, ShouldBeNil)
			})

			Convey("getting the connection string returns the correct hosts", func() {
				expected := "localhost:3030"
				state.EXPECT().CurrentConnectionString().Return([]string{expected})
				res, err := conn.CurrentConnectionString()
				So(err, ShouldBeNil)
				So(res, ShouldEqual, expected)
			})

		})

		Convey("allow adding a parent watcher", func() {
			expected := make(chan<- zk.Event)

			state.EXPECT().AddParentWatcher(expected)
			state.EXPECT().Close().AnyTimes()

			conn.AddParentWatcher(expected)
		})

		Convey("remove a parent watcher", func() {
			expected := make(chan<- zk.Event)

			state.EXPECT().RemoveParentWatcher(expected)
			state.EXPECT().Close().AnyTimes()

			conn.RemoveParentWatcher(expected)
		})

		Convey("forward the ConnectionIndex", func() {
			expected := int32(3039)
			state.EXPECT().InstanceIndex().Return(expected)
			state.EXPECT().Close().AnyTimes()

			So(conn.ConnectionIndex(), ShouldEqual, expected)
		})

		Convey("forward the IsConnected status", func() {
			expected := true
			state.EXPECT().IsConnected().Return(expected)
			state.EXPECT().Close().AnyTimes()

			So(conn.IsConnected(), ShouldBeTrue)
		})
	})

}
