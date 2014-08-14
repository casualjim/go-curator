// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/obeattie/go-zookeeper/zk (interfaces: IConn)

package mock_zk

import (
	zk "github.com/obeattie/go-zookeeper/zk"
	gomock "code.google.com/p/gomock/gomock"
)

// Mock of IConn interface
type MockIConn struct {
	ctrl     *gomock.Controller
	recorder *_MockIConnRecorder
}

// Recorder for MockIConn (not exported)
type _MockIConnRecorder struct {
	mock *MockIConn
}

func NewMockIConn(ctrl *gomock.Controller) *MockIConn {
	mock := &MockIConn{ctrl: ctrl}
	mock.recorder = &_MockIConnRecorder{mock}
	return mock
}

func (_m *MockIConn) EXPECT() *_MockIConnRecorder {
	return _m.recorder
}

func (_m *MockIConn) AddAuth(_param0 string, _param1 []byte) error {
	ret := _m.ctrl.Call(_m, "AddAuth", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIConnRecorder) AddAuth(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddAuth", arg0, arg1)
}

func (_m *MockIConn) Children(_param0 string) ([]string, *zk.Stat, error) {
	ret := _m.ctrl.Call(_m, "Children", _param0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(*zk.Stat)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockIConnRecorder) Children(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Children", arg0)
}

func (_m *MockIConn) ChildrenW(_param0 string) ([]string, *zk.Stat, <-chan zk.Event, error) {
	ret := _m.ctrl.Call(_m, "ChildrenW", _param0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(*zk.Stat)
	ret2, _ := ret[2].(<-chan zk.Event)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

func (_mr *_MockIConnRecorder) ChildrenW(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ChildrenW", arg0)
}

func (_m *MockIConn) Close() {
	_m.ctrl.Call(_m, "Close")
}

func (_mr *_MockIConnRecorder) Close() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Close")
}

func (_m *MockIConn) Create(_param0 string, _param1 []byte, _param2 int32, _param3 []zk.ACL) (string, error) {
	ret := _m.ctrl.Call(_m, "Create", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIConnRecorder) Create(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Create", arg0, arg1, arg2, arg3)
}

func (_m *MockIConn) CreateProtectedEphemeralSequential(_param0 string, _param1 []byte, _param2 []zk.ACL) (string, error) {
	ret := _m.ctrl.Call(_m, "CreateProtectedEphemeralSequential", _param0, _param1, _param2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIConnRecorder) CreateProtectedEphemeralSequential(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateProtectedEphemeralSequential", arg0, arg1, arg2)
}

func (_m *MockIConn) Delete(_param0 string, _param1 int32) error {
	ret := _m.ctrl.Call(_m, "Delete", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIConnRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Delete", arg0, arg1)
}

func (_m *MockIConn) Exists(_param0 string) (bool, *zk.Stat, error) {
	ret := _m.ctrl.Call(_m, "Exists", _param0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*zk.Stat)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockIConnRecorder) Exists(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Exists", arg0)
}

func (_m *MockIConn) ExistsW(_param0 string) (bool, *zk.Stat, <-chan zk.Event, error) {
	ret := _m.ctrl.Call(_m, "ExistsW", _param0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*zk.Stat)
	ret2, _ := ret[2].(<-chan zk.Event)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

func (_mr *_MockIConnRecorder) ExistsW(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ExistsW", arg0)
}

func (_m *MockIConn) Get(_param0 string) ([]byte, *zk.Stat, error) {
	ret := _m.ctrl.Call(_m, "Get", _param0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(*zk.Stat)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockIConnRecorder) Get(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Get", arg0)
}

func (_m *MockIConn) GetACL(_param0 string) ([]zk.ACL, *zk.Stat, error) {
	ret := _m.ctrl.Call(_m, "GetACL", _param0)
	ret0, _ := ret[0].([]zk.ACL)
	ret1, _ := ret[1].(*zk.Stat)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockIConnRecorder) GetACL(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetACL", arg0)
}

func (_m *MockIConn) GetW(_param0 string) ([]byte, *zk.Stat, <-chan zk.Event, error) {
	ret := _m.ctrl.Call(_m, "GetW", _param0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(*zk.Stat)
	ret2, _ := ret[2].(<-chan zk.Event)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

func (_mr *_MockIConnRecorder) GetW(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetW", arg0)
}

func (_m *MockIConn) Multi(_param0 zk.MultiOps) error {
	ret := _m.ctrl.Call(_m, "Multi", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIConnRecorder) Multi(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Multi", arg0)
}

func (_m *MockIConn) Reconnect() error {
	ret := _m.ctrl.Call(_m, "Reconnect")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIConnRecorder) Reconnect() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Reconnect")
}

func (_m *MockIConn) Set(_param0 string, _param1 []byte, _param2 int32) (*zk.Stat, error) {
	ret := _m.ctrl.Call(_m, "Set", _param0, _param1, _param2)
	ret0, _ := ret[0].(*zk.Stat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIConnRecorder) Set(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Set", arg0, arg1, arg2)
}

func (_m *MockIConn) SetACL(_param0 string, _param1 []zk.ACL, _param2 int32) (*zk.Stat, error) {
	ret := _m.ctrl.Call(_m, "SetACL", _param0, _param1, _param2)
	ret0, _ := ret[0].(*zk.Stat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIConnRecorder) SetACL(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetACL", arg0, arg1, arg2)
}

func (_m *MockIConn) State() zk.State {
	ret := _m.ctrl.Call(_m, "State")
	ret0, _ := ret[0].(zk.State)
	return ret0
}

func (_mr *_MockIConnRecorder) State() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "State")
}

func (_m *MockIConn) Sync(_param0 string) (string, error) {
	ret := _m.ctrl.Call(_m, "Sync", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIConnRecorder) Sync(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Sync", arg0)
}
