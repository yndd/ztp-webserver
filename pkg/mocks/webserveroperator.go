// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/webserver/interfaces/webserveroperator.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockWebserverOperations is a mock of WebserverOperations interface.
type MockWebserverOperations struct {
	ctrl     *gomock.Controller
	recorder *MockWebserverOperationsMockRecorder
}

// MockWebserverOperationsMockRecorder is the mock recorder for MockWebserverOperations.
type MockWebserverOperationsMockRecorder struct {
	mock *MockWebserverOperations
}

// NewMockWebserverOperations creates a new mock instance.
func NewMockWebserverOperations(ctrl *gomock.Controller) *MockWebserverOperations {
	mock := &MockWebserverOperations{ctrl: ctrl}
	mock.recorder = &MockWebserverOperationsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebserverOperations) EXPECT() *MockWebserverOperationsMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockWebserverOperations) Run(port int, storageFolder string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", port, storageFolder)
}

// Run indicates an expected call of Run.
func (mr *MockWebserverOperationsMockRecorder) Run(port, storageFolder interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockWebserverOperations)(nil).Run), port, storageFolder)
}

// SetKubeConfig mocks base method.
func (m *MockWebserverOperations) SetKubeConfig(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetKubeConfig", arg0)
}

// SetKubeConfig indicates an expected call of SetKubeConfig.
func (mr *MockWebserverOperationsMockRecorder) SetKubeConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetKubeConfig", reflect.TypeOf((*MockWebserverOperations)(nil).SetKubeConfig), arg0)
}
