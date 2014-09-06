package events

import (
	"testing"
	"time"

	"github.com/obeattie/go-zookeeper/zk"
	//. "github.com/onsi/gomega"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEventBus(t *testing.T) {
	Convey("An eventbus should", t, func() {
		bus := New()

		Reset(func() {
			bus.Close()
		})

		Convey("register handlers", func() {
			h := make(chan zk.Event)
			So(bus.Len(), ShouldEqual, 0)
			bus.Add(h)
			So(bus.Len(), ShouldEqual, 1)
		})

		Convey("unregister handlers", func() {
			h := make(chan zk.Event)
			h2 := make(chan zk.Event)
			h3 := make(chan zk.Event)
			So(bus.Len(), ShouldEqual, 0)
			bus.Add(h, h2, h3)
			So(bus.Len(), ShouldEqual, 3)
			bus.Remove(h2)
			So(bus.Len(), ShouldEqual, 2)
		})

		Convey("when publishing events", func() {
			listener1 := make(chan zk.Event)
			listener2 := make(chan zk.Event)
			listener3 := make(chan zk.Event)
			bus.Add(listener1, listener2, listener3)

			Convey("send events to all registered listeners", func() {
				evt := zk.Event{
					Type:  zk.EventSession,
					State: zk.StateHasSession,
					Path:  "",
					Err:   nil,
				}

				evts := make([]zk.Event, 3)
				seen := 0
				latch := make(chan bool)
				go func() {
					for seen < 3 {
						select {
						case e1 := <-listener1:
							evts[0] = e1
							seen++
						case e2 := <-listener2:
							evts[1] = e2
							seen++
						case e3 := <-listener3:
							evts[2] = e3
							seen++
						}
					}
					latch <- true
				}()

				bus.Trigger() <- evt
				<-latch
				So(evts[0], ShouldResemble, evt)
				So(evts[1], ShouldResemble, evt)
				So(evts[2], ShouldResemble, evt)
			})

			Convey("drop listeners when they don't accept the message", func() {
				evt := zk.Event{
					Type:  zk.EventSession,
					State: zk.StateHasSession,
					Path:  "",
					Err:   nil,
				}

				latch := make(chan bool)
				go func() {
					time.Sleep(1 * time.Second)
					latch <- true
				}()

				bus.Trigger() <- evt
				<-latch

				So(bus.Len(), ShouldEqual, 0)
			})
		})
	})
}
