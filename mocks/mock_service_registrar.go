// Code generated by MockGen. DO NOT EDIT.
// Source: google.golang.org/grpc (interfaces: ServiceRegistrar)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockServiceRegistrar is a mock of ServiceRegistrar interface.
type MockServiceRegistrar struct {
	ctrl     *gomock.Controller
	recorder *MockServiceRegistrarMockRecorder
}

// MockServiceRegistrarMockRecorder is the mock recorder for MockServiceRegistrar.
type MockServiceRegistrarMockRecorder struct {
	mock *MockServiceRegistrar
}

// NewMockServiceRegistrar creates a new mock instance.
func NewMockServiceRegistrar(ctrl *gomock.Controller) *MockServiceRegistrar {
	mock := &MockServiceRegistrar{ctrl: ctrl}
	mock.recorder = &MockServiceRegistrarMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceRegistrar) EXPECT() *MockServiceRegistrarMockRecorder {
	return m.recorder
}

// RegisterService mocks base method.
func (m *MockServiceRegistrar) RegisterService(arg0 *grpc.ServiceDesc, arg1 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterService", arg0, arg1)
}

// RegisterService indicates an expected call of RegisterService.
func (mr *MockServiceRegistrarMockRecorder) RegisterService(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterService", reflect.TypeOf((*MockServiceRegistrar)(nil).RegisterService), arg0, arg1)
}
