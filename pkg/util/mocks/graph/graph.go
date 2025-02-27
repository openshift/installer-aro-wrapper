// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openshift/installer-aro-wrapper/pkg/cluster/graph (interfaces: Manager)

// Package mock_graph is a generated GoMock package.
package mock_graph

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	graph "github.com/openshift/installer-aro-wrapper/pkg/cluster/graph"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Exists mocks base method.
func (m *MockManager) Exists(arg0 context.Context, arg1, arg2 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockManagerMockRecorder) Exists(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockManager)(nil).Exists), arg0, arg1, arg2)
}

// GetUserDelegatedSASIgnitionBlobURL mocks base method.
func (m *MockManager) GetUserDelegatedSASIgnitionBlobURL(arg0 context.Context, arg1, arg2, arg3 string, arg4 bool) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserDelegatedSASIgnitionBlobURL", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserDelegatedSASIgnitionBlobURL indicates an expected call of GetUserDelegatedSASIgnitionBlobURL.
func (mr *MockManagerMockRecorder) GetUserDelegatedSASIgnitionBlobURL(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserDelegatedSASIgnitionBlobURL", reflect.TypeOf((*MockManager)(nil).GetUserDelegatedSASIgnitionBlobURL), arg0, arg1, arg2, arg3, arg4)
}

// LoadPersisted mocks base method.
func (m *MockManager) LoadPersisted(arg0 context.Context, arg1, arg2 string) (graph.PersistedGraph, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadPersisted", arg0, arg1, arg2)
	ret0, _ := ret[0].(graph.PersistedGraph)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadPersisted indicates an expected call of LoadPersisted.
func (mr *MockManagerMockRecorder) LoadPersisted(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadPersisted", reflect.TypeOf((*MockManager)(nil).LoadPersisted), arg0, arg1, arg2)
}

// Save mocks base method.
func (m *MockManager) Save(arg0 context.Context, arg1, arg2 string, arg3 graph.Graph) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockManagerMockRecorder) Save(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockManager)(nil).Save), arg0, arg1, arg2, arg3)
}
