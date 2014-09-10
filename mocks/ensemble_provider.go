// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/casualjim/go-curator/ensemble (interfaces: Provider)

package mocks

import (
	gomock "code.google.com/p/gomock/gomock"
)

// Mock of Provider interface
type MockEnsembleProvider struct {
	ctrl     *gomock.Controller
	recorder *_MockEnsembleProviderRecorder
}

// Recorder for MockEnsembleProvider (not exported)
type _MockEnsembleProviderRecorder struct {
	mock *MockEnsembleProvider
}

func NewMockEnsembleProvider(ctrl *gomock.Controller) *MockEnsembleProvider {
	mock := &MockEnsembleProvider{ctrl: ctrl}
	mock.recorder = &_MockEnsembleProviderRecorder{mock}
	return mock
}

func (_m *MockEnsembleProvider) EXPECT() *_MockEnsembleProviderRecorder {
	return _m.recorder
}

func (_m *MockEnsembleProvider) Close() error {
	ret := _m.ctrl.Call(_m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockEnsembleProviderRecorder) Close() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Close")
}

func (_m *MockEnsembleProvider) Hosts() []string {
	ret := _m.ctrl.Call(_m, "Hosts")
	ret0, _ := ret[0].([]string)
	return ret0
}

func (_mr *_MockEnsembleProviderRecorder) Hosts() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Hosts")
}

func (_m *MockEnsembleProvider) Start() error {
	ret := _m.ctrl.Call(_m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockEnsembleProviderRecorder) Start() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Start")
}
