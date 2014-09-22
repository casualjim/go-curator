package shared

import "github.com/casualjim/go-zookeeper/zk"

// WatcherHolder is a construct to allow gomock to ignore calls
type WatcherHolder struct {
	Watcher chan zk.Event
}
