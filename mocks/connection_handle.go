// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/casualjim/go-curator (interfaces: ConnectionHandle)

package mocks

import (
	zk "github.com/casualjim/go-zookeeper/zk"
	time "time"
	gomock "code.google.com/p/gomock/gomock"
)

// Mock of ConnectionHandle interface
type MockConnectionHandle struct {
	ctrl     *gomock.Controller
	recorder *_MockConnectionHandleRecorder
}

// Recorder for MockConnectionHandle (not exported)
type _MockConnectionHandleRecorder struct {
	mock *MockConnectionHandle
}

func NewMockConnectionHandle(ctrl *gomock.Controller) *MockConnectionHandle {
	mock := &MockConnectionHandle{ctrl: ctrl}
	mock.recorder = &_MockConnectionHandleRecorder{mock}
	return mock
}

func (_m *MockConnectionHandle) EXPECT() *_MockConnectionHandleRecorder {
	return _m.recorder
}

func (_m *MockConnectionHandle) Close() error {
	ret := _m.ctrl.Call(_m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockConnectionHandleRecorder) Close() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Close")
}

func (_m *MockConnectionHandle) Conn() zk.IConn {
	ret := _m.ctrl.Call(_m, "Conn")
	ret0, _ := ret[0].(zk.IConn)
	return ret0
}

func (_mr *_MockConnectionHandleRecorder) Conn() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Conn")
}

func (_m *MockConnectionHandle) HasNewConnectionString() bool {
	ret := _m.ctrl.Call(_m, "HasNewConnectionString")
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockConnectionHandleRecorder) HasNewConnectionString() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "HasNewConnectionString")
}

func (_m *MockConnectionHandle) Hosts() []string {
	ret := _m.ctrl.Call(_m, "Hosts")
	ret0, _ := ret[0].([]string)
	return ret0
}

func (_mr *_MockConnectionHandleRecorder) Hosts() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Hosts")
}

func (_m *MockConnectionHandle) Reconnect() (<-chan zk.Event, error) {
	ret := _m.ctrl.Call(_m, "Reconnect")
	ret0, _ := ret[0].(<-chan zk.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockConnectionHandleRecorder) Reconnect() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Reconnect")
}

func (_m *MockConnectionHandle) SessionTimeout() time.Duration {
	ret := _m.ctrl.Call(_m, "SessionTimeout")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

func (_mr *_MockConnectionHandleRecorder) SessionTimeout() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SessionTimeout")
}
