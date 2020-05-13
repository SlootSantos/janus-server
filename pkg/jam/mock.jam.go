// Code generated by MockGen. DO NOT EDIT.
// Source: jam.go

// Package jam is a generated GoMock package.
package jam

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockstackResource is a mock of stackResource interface
type MockstackResource struct {
	ctrl     *gomock.Controller
	recorder *MockstackResourceMockRecorder
}

// MockstackResourceMockRecorder is the mock recorder for MockstackResource
type MockstackResourceMockRecorder struct {
	mock *MockstackResource
}

// NewMockstackResource creates a new mock instance
func NewMockstackResource(ctrl *gomock.Controller) *MockstackResource {
	mock := &MockstackResource{ctrl: ctrl}
	mock.recorder = &MockstackResourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockstackResource) EXPECT() *MockstackResourceMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockstackResource) Create(arg0 context.Context, arg1 *CreationParam, arg2 *OutputParam) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockstackResourceMockRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockstackResource)(nil).Create), arg0, arg1, arg2)
}

// Destroy mocks base method
func (m *MockstackResource) Destroy(arg0 context.Context, arg1 *DeletionParam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy
func (mr *MockstackResourceMockRecorder) Destroy(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockstackResource)(nil).Destroy), arg0, arg1)
}

// List mocks base method
func (m *MockstackResource) List(arg0 context.Context) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// List indicates an expected call of List
func (mr *MockstackResourceMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockstackResource)(nil).List), arg0)
}
