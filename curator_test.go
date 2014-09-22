package curator

import (
	"errors"
	"testing"
	"time"

	"code.google.com/p/gomock/gomock"

	"github.com/casualjim/go-curator/mocks"
	"github.com/casualjim/go-zookeeper/zk"
	. "github.com/smartystreets/goconvey/convey"
)

type recyclingWatcherFactory struct {
	current chan zk.Event
}

func (r *recyclingWatcherFactory) MakeWatcher() chan zk.Event {
	if r.current == nil {
		r.current = make(chan zk.Event)
	}
	return r.current
}

func (r *recyclingWatcherFactory) Invalidate() {
	if r.current != nil {
		close(r.current)
		r.current = nil
	}
}

func TestCurator(t *testing.T) {

	Convey("The curator zookeeper client should", t, func() {
		//RegisterTestingT(t)
		mockCtrl := gomock.NewController(t)
		state := mocks.NewMockConnectionState(mockCtrl)
		conn := &Conn{
			state:             state,
			RetryPolicy:       ForeverPolicy,
			ConnectionTimeout: 3 * time.Second,
			started:           0,
			watcherFactory:    &recyclingWatcherFactory{},
		}
		conn.startedPtr = &conn.started

		Reset(func() {
			mockCtrl.Finish()
			conn.Close()
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

			Convey("BlockUntilConnectedOrTimedOut returns immediately when connected", func() {
				state.EXPECT().IsConnected().Times(2).Return(true)
				res, err := conn.BlockUntilConnectedOrTimedOut()
				So(err, ShouldBeNil)
				So(res, ShouldBeTrue)
			})

			Convey("BlockUntilConnectedOrTimedOut eventually returns true", func() {

				gomock.InOrder(
					state.EXPECT().IsConnected().Return(false),
					state.EXPECT().IsConnected().Return(false),
					state.EXPECT().IsConnected().Return(true).AnyTimes(),
				)
				state.EXPECT().AddParentWatcherHolder(gomock.Any()).AnyTimes()
				// state.EXPECT().IsConnected().Return(false)
				// state.EXPECT().IsConnected().Return(true)
				// state.EXPECT().IsConnected().Return(true)

				result := make(chan error)
				go func() {
					res, err := conn.BlockUntilConnectedOrTimedOut()
					logger.Debug("Returned with: %v, %v", res, err)
					if !res {
						result <- errors.New("Not connected!")
					} else {
						if err != nil {
							result <- err
						} else {
							result <- nil
						}
					}
				}()

				err2 := make(chan error)
				go func() {
					for {
						select {
						case r := <-result:
							err2 <- r
							break
						case <-time.After(10 * time.Second):
							err2 <- errors.New("Timed out!")
							break
						}
					}
				}()

				actual := <-err2
				So(actual, ShouldBeNil)

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
