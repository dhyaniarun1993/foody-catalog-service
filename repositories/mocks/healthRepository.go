// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dhyaniarun1993/foody-catalog-service/repositories (interfaces: HealthRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	errors "github.com/dhyaniarun1993/foody-common/errors"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockHealthRepository is a mock of HealthRepository interface
type MockHealthRepository struct {
	ctrl     *gomock.Controller
	recorder *MockHealthRepositoryMockRecorder
}

// MockHealthRepositoryMockRecorder is the mock recorder for MockHealthRepository
type MockHealthRepositoryMockRecorder struct {
	mock *MockHealthRepository
}

// NewMockHealthRepository creates a new mock instance
func NewMockHealthRepository(ctrl *gomock.Controller) *MockHealthRepository {
	mock := &MockHealthRepository{ctrl: ctrl}
	mock.recorder = &MockHealthRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHealthRepository) EXPECT() *MockHealthRepositoryMockRecorder {
	return m.recorder
}

// HealthCheck mocks base method
func (m *MockHealthRepository) HealthCheck(arg0 context.Context) errors.AppError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthCheck", arg0)
	ret0, _ := ret[0].(errors.AppError)
	return ret0
}

// HealthCheck indicates an expected call of HealthCheck
func (mr *MockHealthRepositoryMockRecorder) HealthCheck(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthCheck", reflect.TypeOf((*MockHealthRepository)(nil).HealthCheck), arg0)
}