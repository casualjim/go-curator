package events

import (
	"io"
	"sync"
	"time"

	"github.com/casualjim/go-zookeeper/zk"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("events")

// EventBus does fanout to registered channels
type EventBus interface {
	io.Closer
	Trigger() chan<- zk.Event
	Add(handlers ...chan<- zk.Event)
	Remove(handlers ...chan<- zk.Event)
	Len() int
}

type defaultEventBus struct {
	sync.Mutex

	channel  chan zk.Event
	handlers []chan<- zk.Event
	closing  chan chan struct{}
}

func New() EventBus { return NewWithTimeout(50 * time.Millisecond) }

// NewWithTimeout creates a new event bus with a custom timeout for message delivery
func NewWithTimeout(timeout time.Duration) EventBus {
	e := &defaultEventBus{
		closing:  make(chan chan struct{}),
		channel:  make(chan zk.Event),
		handlers: []chan<- zk.Event{},
	}
	go e.dispatcherLoop(timeout)
	return e
}

func (e *defaultEventBus) dispatcherLoop(timeout time.Duration) {
	for {
		select {
		case evt := <-e.channel:
			for _, handler := range e.handlers {
				go e.dispatchEventWithTimeout(handler, timeout, evt)
			}
		case closed := <-e.closing:
			close(e.channel)
			for _, handler := range e.handlers {
				close(handler)
			}
			e.handlers = []chan<- zk.Event{}
			closed <- struct{}{}
			return
		}
	}
}

func (e *defaultEventBus) dispatchEventWithTimeout(channel chan<- zk.Event, timeout time.Duration, event zk.Event) {
	timer := time.NewTimer(timeout)
	select {
	case channel <- event:
		timer.Stop()
	case <-timer.C:
		log.Warning("Failed to send event %+v to listener within %v", event, timeout)
		e.Remove(channel)
	}
}

func (e *defaultEventBus) Trigger() chan<- zk.Event {
	return e.channel
}

func (e *defaultEventBus) Add(handler ...chan<- zk.Event) {
	e.Lock()
	defer e.Unlock()
	e.handlers = append(e.handlers, handler...)
}

func (e *defaultEventBus) Remove(handler ...chan<- zk.Event) {
	e.Lock()
	defer e.Unlock()
	for _, h := range handler {
		for i, handler := range e.handlers {
			if h == handler {
				e.handlers = append(e.handlers[:i], e.handlers[i+1:]...)
				break
			}
		}
	}
}

func (e *defaultEventBus) Close() error {
	ch := make(chan struct{})
	e.closing <- ch
	<-ch
	close(e.closing)

	return nil
}

func (e *defaultEventBus) Len() int {
	return len(e.handlers)
}
