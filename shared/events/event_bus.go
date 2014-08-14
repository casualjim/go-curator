package events

import (
	"io"
	"sync"

	"github.com/obeattie/go-zookeeper/zk"
)

// EventBus does fanout to registered channels
type EventBus interface {
	io.Closer
	Trigger() chan<- zk.Event
	Add(handlers ...chan<- zk.Event)
	Remove(handlers ...chan<- zk.Event)
}

type defaultEventBus struct {
	sync.Mutex

	channel  chan zk.Event
	handlers []chan<- zk.Event
	closing  chan chan struct{}
}

// New creates a new event bus
func New() EventBus {
	e := &defaultEventBus{
		closing:  make(chan chan struct{}),
		channel:  make(chan zk.Event),
		handlers: []chan<- zk.Event{},
	}
	go func() {
		for {
			select {
			case evt := <-e.channel:
				for _, handler := range e.handlers {
					handler <- evt
				}
			case closed := <-e.closing:
				close(e.channel)
				closed <- struct{}{}
				return
			}
		}
	}()
	return nil
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
	return nil
}
