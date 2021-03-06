// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/webserver/interfaces/webserversetupper.go

// Package mocks is a generated GoMock package.
package mocks

import (
	http "net/http"
	url "net/url"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	structs "github.com/yndd/ztp-dhcp/pkg/structs"
	interfaces "github.com/yndd/ztp-webserver/pkg/storage/interfaces"
	structs0 "github.com/yndd/ztp-webserver/pkg/structs"
)

// MockWebserverSetupper is a mock of WebserverSetupper interface.
type MockWebserverSetupper struct {
	ctrl     *gomock.Controller
	recorder *MockWebserverSetupperMockRecorder
}

// MockWebserverSetupperMockRecorder is the mock recorder for MockWebserverSetupper.
type MockWebserverSetupperMockRecorder struct {
	mock *MockWebserverSetupper
}

// NewMockWebserverSetupper creates a new mock instance.
func NewMockWebserverSetupper(ctrl *gomock.Controller) *MockWebserverSetupper {
	mock := &MockWebserverSetupper{ctrl: ctrl}
	mock.recorder = &MockWebserverSetupperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebserverSetupper) EXPECT() *MockWebserverSetupperMockRecorder {
	return m.recorder
}

// AddHandler mocks base method.
func (m *MockWebserverSetupper) AddHandler(arg0 *structs0.UrlParams, arg1 func(http.ResponseWriter, *http.Request)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddHandler", arg0, arg1)
}

// AddHandler indicates an expected call of AddHandler.
func (mr *MockWebserverSetupperMockRecorder) AddHandler(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddHandler", reflect.TypeOf((*MockWebserverSetupper)(nil).AddHandler), arg0, arg1)
}

// EnrichUrl mocks base method.
func (m *MockWebserverSetupper) EnrichUrl(arg0 *url.URL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnrichUrl", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnrichUrl indicates an expected call of EnrichUrl.
func (mr *MockWebserverSetupperMockRecorder) EnrichUrl(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnrichUrl", reflect.TypeOf((*MockWebserverSetupper)(nil).EnrichUrl), arg0)
}

// GetDeviceInformationByName mocks base method.
func (m *MockWebserverSetupper) GetDeviceInformationByName(deviceId string) (*structs.DeviceInformation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceInformationByName", deviceId)
	ret0, _ := ret[0].(*structs.DeviceInformation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceInformationByName indicates an expected call of GetDeviceInformationByName.
func (mr *MockWebserverSetupperMockRecorder) GetDeviceInformationByName(deviceId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceInformationByName", reflect.TypeOf((*MockWebserverSetupper)(nil).GetDeviceInformationByName), deviceId)
}

// GetIndex mocks base method.
func (m *MockWebserverSetupper) GetIndex() interfaces.Index {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIndex")
	ret0, _ := ret[0].(interfaces.Index)
	return ret0
}

// GetIndex indicates an expected call of GetIndex.
func (mr *MockWebserverSetupperMockRecorder) GetIndex() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIndex", reflect.TypeOf((*MockWebserverSetupper)(nil).GetIndex))
}

// GetStorage mocks base method.
func (m *MockWebserverSetupper) GetStorage() interfaces.Storage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStorage")
	ret0, _ := ret[0].(interfaces.Storage)
	return ret0
}

// GetStorage indicates an expected call of GetStorage.
func (mr *MockWebserverSetupperMockRecorder) GetStorage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStorage", reflect.TypeOf((*MockWebserverSetupper)(nil).GetStorage))
}

// ResponseFromIndex mocks base method.
func (m *MockWebserverSetupper) ResponseFromIndex(arg0 http.ResponseWriter, arg1 *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResponseFromIndex", arg0, arg1)
}

// ResponseFromIndex indicates an expected call of ResponseFromIndex.
func (mr *MockWebserverSetupperMockRecorder) ResponseFromIndex(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResponseFromIndex", reflect.TypeOf((*MockWebserverSetupper)(nil).ResponseFromIndex), arg0, arg1)
}
