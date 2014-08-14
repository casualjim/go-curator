package main

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func main() {
	conn, ch, err := zk.Connect([]string{"exeggutor-box:2181"}, 5*time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	closeCh := make(chan struct{})
	var i int
	go func() {
		for {
			select {
			case ev := <-ch:
				fmt.Printf("Got event %+v\n", ev)
				if ev.State == zk.StateDisconnected {
					if i > 30 {
						closeCh <- struct{}{}
					}
					i++
				} else {
					i = 0
				}
			case <-time.After(5 * time.Second):
				fmt.Println("Closing connection")
				conn.Close()
			}
		}
	}()

	<-closeCh

}
