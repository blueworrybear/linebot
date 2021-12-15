// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/blueworrybear/lineBot/core (interfaces: UserStore)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	core "github.com/blueworrybear/lineBot/core"
	gomock "github.com/golang/mock/gomock"
)

// MockUserStore is a mock of UserStore interface.
type MockUserStore struct {
	ctrl     *gomock.Controller
	recorder *MockUserStoreMockRecorder
}

// MockUserStoreMockRecorder is the mock recorder for MockUserStore.
type MockUserStoreMockRecorder struct {
	mock *MockUserStore
}

// NewMockUserStore creates a new mock instance.
func NewMockUserStore(ctrl *gomock.Controller) *MockUserStore {
	mock := &MockUserStore{ctrl: ctrl}
	mock.recorder = &MockUserStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserStore) EXPECT() *MockUserStoreMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockUserStore) All() ([]*core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All")
	ret0, _ := ret[0].([]*core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockUserStoreMockRecorder) All() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockUserStore)(nil).All))
}

// Create mocks base method.
func (m *MockUserStore) Create(arg0 *core.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUserStoreMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserStore)(nil).Create), arg0)
}

// Find mocks base method.
func (m *MockUserStore) Find(arg0 *core.User) (*core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0)
	ret0, _ := ret[0].(*core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockUserStoreMockRecorder) Find(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockUserStore)(nil).Find), arg0)
}

// FindForRequest mocks base method.
func (m *MockUserStore) FindForRequest(arg0 *core.User) (*core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindForRequest", arg0)
	ret0, _ := ret[0].(*core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindForRequest indicates an expected call of FindForRequest.
func (mr *MockUserStoreMockRecorder) FindForRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindForRequest", reflect.TypeOf((*MockUserStore)(nil).FindForRequest), arg0)
}

// Update mocks base method.
func (m *MockUserStore) Update(arg0 *core.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserStoreMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserStore)(nil).Update), arg0)
}

// VIPs mocks base method.
func (m *MockUserStore) VIPs() ([]*core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VIPs")
	ret0, _ := ret[0].([]*core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VIPs indicates an expected call of VIPs.
func (mr *MockUserStoreMockRecorder) VIPs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VIPs", reflect.TypeOf((*MockUserStore)(nil).VIPs))
}
