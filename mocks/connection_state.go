// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/casualjim/go-curator (interfaces: ConnectionState)

package mocks

import (
	gomock "code.google.com/p/gomock/gomock"
	"github.com/casualjim/go-curator/shared"
	zk "github.com/casualjim/go-zookeeper/zk"
)

// Mock of ConnectionState interface
type MockConnectionState struct {
	ctrl     *gomock.Controller
	recorder *_MockConnectionStateRecorder
}

// Recorder for MockConnectionState (not exported)
type _MockConnectionStateRecorder struct {
	mock *MockConnectionState
}

func NewMockConnectionState(ctrl *gomock.Controller) *MockConnectionState {
	mock := &MockConnectionState{ctrl: ctrl}
	mock.recorder = &_MockConnectionStateRecorder{mock}
	return mock
}

func (_m *MockConnectionState) EXPECT() *_MockConnectionStateRecorder {
	return _m.recorder
}

func (_m *MockConnectionState) AddParentWatcher(_param0 chan<- zk.Event) {
	_m.ctrl.Call(_m, "AddParentWatcher", _param0)
}

func (_mr *_MockConnectionStateRecorder) AddParentWatcher(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddParentWatcher", arg0)
}

func (_m *MockConnectionState) AddParentWatcherHolder(_param0 *shared.WatcherHolder) {
	_m.ctrl.Call(_m, "AddParentWatcherHolder", _param0)
}

func (_mr *_MockConnectionStateRecorder) AddParentWatcherHolder(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddParentWatcherHolder", arg0)
}

func (_m *MockConnectionState) Close() error {
	ret := _m.ctrl.Call(_m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockConnectionStateRecorder) Close() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Close")
}

func (_m *MockConnectionState) Conn() (zk.IConn, error) {
	ret := _m.ctrl.Call(_m, "Conn")
	ret0, _ := ret[0].(zk.IConn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockConnectionStateRecorder) Conn() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Conn")
}

func (_m *MockConnectionState) CurrentConnectionString() []string {
	ret := _m.ctrl.Call(_m, "CurrentConnectionString")
	ret0, _ := ret[0].([]string)
	return ret0
}

func (_mr *_MockConnectionStateRecorder) CurrentConnectionString() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CurrentConnectionString")
}

func (_m *MockConnectionState) InstanceIndex() int32 {
	ret := _m.ctrl.Call(_m, "InstanceIndex")
	ret0, _ := ret[0].(int32)
	return ret0
}

func (_mr *_MockConnectionStateRecorder) InstanceIndex() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "InstanceIndex")
}

func (_m *MockConnectionState) IsConnected() bool {
	ret := _m.ctrl.Call(_m, "IsConnected")
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockConnectionStateRecorder) IsConnected() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "IsConnected")
}

func (_m *MockConnectionState) RemoveParentWatcher(_param0 chan<- zk.Event) {
	_m.ctrl.Call(_m, "RemoveParentWatcher", _param0)
}

func (_mr *_MockConnectionStateRecorder) RemoveParentWatcher(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemoveParentWatcher", arg0)
}

func (_m *MockConnectionState) RemoveParentWatcherHolder(_param0 *shared.WatcherHolder) {
	_m.ctrl.Call(_m, "RemoveParentWatcherHolder", _param0)
}

func (_mr *_MockConnectionStateRecorder) RemoveParentWatcherHolder(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemoveParentWatcherHolder", arg0)
}

func (_m *MockConnectionState) Start() error {
	ret := _m.ctrl.Call(_m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockConnectionStateRecorder) Start() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Start")
}
