// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/storage/interfaces/storage.go

// Package mocks is a generated GoMock package.
package mocks

import (
	fs "io/fs"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockStorage) Handle(w http.ResponseWriter, filePath string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Handle", w, filePath)
}

// Handle indicates an expected call of Handle.
func (mr *MockStorageMockRecorder) Handle(w, filePath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockStorage)(nil).Handle), w, filePath)
}

// LoadBackend mocks base method.
func (m *MockStorage) LoadBackend(arg0 fs.FS) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadBackend", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadBackend indicates an expected call of LoadBackend.
func (mr *MockStorageMockRecorder) LoadBackend(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadBackend", reflect.TypeOf((*MockStorage)(nil).LoadBackend), arg0)
}