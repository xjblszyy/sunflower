// Code generated by MockGen. DO NOT EDIT.
// Source: grade.go

// Package mock_dao is a generated GoMock package.
package mock_dao

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	dao "sunflower/pkg/app/apiserver/dao"
	model "sunflower/pkg/app/apiserver/model"
	gormutil_v2 "sunflower/pkg/libs/gormutil_v2"
)

// MockGradeDao is a mock of GradeDao interface
type MockGradeDao struct {
	ctrl     *gomock.Controller
	recorder *MockGradeDaoMockRecorder
}

// MockGradeDaoMockRecorder is the mock recorder for MockGradeDao
type MockGradeDaoMockRecorder struct {
	mock *MockGradeDao
}

// NewMockGradeDao creates a new mock instance
func NewMockGradeDao(ctrl *gomock.Controller) *MockGradeDao {
	mock := &MockGradeDao{ctrl: ctrl}
	mock.recorder = &MockGradeDaoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGradeDao) EXPECT() *MockGradeDaoMockRecorder {
	return m.recorder
}

// List mocks base method
func (m *MockGradeDao) List(scopes *dao.GradeScope, limit, cursor int) ([]model.Grade, *gormutil_v2.Page, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", scopes, limit, cursor)
	ret0, _ := ret[0].([]model.Grade)
	ret1, _ := ret[1].(*gormutil_v2.Page)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List
func (mr *MockGradeDaoMockRecorder) List(scopes, limit, cursor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockGradeDao)(nil).List), scopes, limit, cursor)
}

// CreateMany mocks base method
func (m *MockGradeDao) CreateMany(grads []model.Grade) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMany", grads)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateMany indicates an expected call of CreateMany
func (mr *MockGradeDaoMockRecorder) CreateMany(grads interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMany", reflect.TypeOf((*MockGradeDao)(nil).CreateMany), grads)
}

// FetchOne mocks base method
func (m *MockGradeDao) FetchOne(scopes *dao.GradeScope) (model.Grade, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchOne", scopes)
	ret0, _ := ret[0].(model.Grade)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchOne indicates an expected call of FetchOne
func (mr *MockGradeDaoMockRecorder) FetchOne(scopes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchOne", reflect.TypeOf((*MockGradeDao)(nil).FetchOne), scopes)
}
