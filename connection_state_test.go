package curator

import (
	"fmt"
	"testing"
	"time"

	"code.google.com/p/gomock/gomock"

	"github.com/casualjim/go-curator/ensemble"
	"github.com/casualjim/go-zookeeper/zk"
	. "github.com/onsi/gomega"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConnectionState(t *testing.T) {

	Convey("The connection state should", t, func() {
		RegisterTestingT(t)
		mockCtrl := gomock.NewController(t)
		factory := testFactory(mockCtrl)
		prov, _, _ := ensemble.Fixed("zk://localhost:39449/blah")
		state := newConnectionState(factory, prov, 1*time.Second, 2*time.Second)

		Reset(func() {
			state.Close()
			mockCtrl.Finish()
		})

		Convey("when starting", func() {
			state.Start()
			factory.conn.EXPECT().Close().AnyTimes()

			Convey("it increments the count", func() {
				So(state.InstanceIndex(), ShouldEqual, 1)
			})

			Convey("it connects to zookeeper", func() {
				So(factory.callCount, ShouldEqual, 1)
			})

			Convey("it signals connected when event received", func() {
				go func() {
					factory.channel <- zk.Event{
						Type:  zk.EventSession,
						State: zk.StateSyncConnected,
					}
				}()

				Eventually(func() bool { return state.IsConnected() }).Should(BeTrue())
				So(state.IsConnected(), ShouldBeTrue)
			})

			Convey("it has the current connection string", func() {
				hosts := state.CurrentConnectionString()
				So(len(hosts), ShouldEqual, 1)
				So(hosts[0], ShouldEqual, "localhost:39449")
			})
		})

		Convey("when connected and acquiring a connection", func() {
			state.Start()
			factory.conn.EXPECT().Close().AnyTimes()

			Convey("it returns an error when there is an error on the queue", func() {
				state.errorQueue.PushBack(fmt.Errorf("expected"))
				c, e := state.Conn()
				So(c, ShouldBeNil)
				So(e, ShouldNotBeNil)
			})

			Convey("it returns an error when there is a timeout, but it's in the tolerance window", func() {
				done := make(chan bool)
				go func() {
					factory.channel <- zk.Event{
						Type:  zk.EventSession,
						State: zk.StateSyncConnected,
					}
					for i := 0; i < 5; i++ {
						if state.IsConnected() {
							done <- true
							return
						}
						<-time.After(500 * time.Millisecond)
					}
					done <- false
				}()
				res := <-done
				So(res, ShouldBeTrue)
				state.started = time.Now().Add(-1500 * time.Millisecond)
				c, e := state.Conn()
				So(c, ShouldBeNil)
				So(e, ShouldNotBeNil)
			})

			Convey("it returns a connection after resetting when there is a timeout greater than the tolerance window", func() {
				done := make(chan bool)
				go func() {
					factory.channel <- zk.Event{
						Type:  zk.EventSession,
						State: zk.StateSyncConnected,
					}
					for i := 0; i < 5; i++ {
						if state.IsConnected() {
							done <- true
							return
						}
						<-time.After(500 * time.Millisecond)
					}
					done <- false
				}()
				res := <-done
				So(res, ShouldBeTrue)

				state.started = time.Now().Add(-2500 * time.Millisecond)
				c, e := state.Conn()
				So(c, ShouldNotBeNil)
				So(e, ShouldBeNil)
				So(state.InstanceIndex(), ShouldEqual, 2)
			})

			Convey("it returns a connection after performing a background reset", func() {
				done := make(chan bool)
				go func() {
					factory.channel <- zk.Event{
						Type:  zk.EventSession,
						State: zk.StateSyncConnected,
					}
					for i := 0; i < 5; i++ {
						if state.IsConnected() {
							done <- true
							return
						}
						<-time.After(500 * time.Millisecond)
					}
					done <- false
				}()
				res := <-done
				So(res, ShouldBeTrue)

				state.started = time.Now().Add(-1500 * time.Millisecond)
				p, _, _ := ensemble.Fixed("zk://localhost:3939/foo")
				state.conn.(*connHandle).ensembleProvider = p
				c, e := state.Conn()
				So(c, ShouldNotBeNil)
				So(e, ShouldBeNil)
				So(state.InstanceIndex(), ShouldEqual, 2)
			})
		})

	})

}
