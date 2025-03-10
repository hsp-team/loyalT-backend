// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go
//
// Generated by this command:
//
//	mockgen -source=handler.go -destination=mocks/mock_service.go -package=mocks coinProgramParticipantService
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	dto "loyalit/internal/domain/entity/dto"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockcoinProgramParticipantService is a mock of coinProgramParticipantService interface.
type MockcoinProgramParticipantService struct {
	ctrl     *gomock.Controller
	recorder *MockcoinProgramParticipantServiceMockRecorder
	isgomock struct{}
}

// MockcoinProgramParticipantServiceMockRecorder is the mock recorder for MockcoinProgramParticipantService.
type MockcoinProgramParticipantServiceMockRecorder struct {
	mock *MockcoinProgramParticipantService
}

// NewMockcoinProgramParticipantService creates a new mock instance.
func NewMockcoinProgramParticipantService(ctrl *gomock.Controller) *MockcoinProgramParticipantService {
	mock := &MockcoinProgramParticipantService{ctrl: ctrl}
	mock.recorder = &MockcoinProgramParticipantServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcoinProgramParticipantService) EXPECT() *MockcoinProgramParticipantServiceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockcoinProgramParticipantService) Get(ctx context.Context, coinProgramParticipantID, userID uuid.UUID) (*dto.CoinProgramParticipantReturn, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, coinProgramParticipantID, userID)
	ret0, _ := ret[0].(*dto.CoinProgramParticipantReturn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockcoinProgramParticipantServiceMockRecorder) Get(ctx, coinProgramParticipantID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockcoinProgramParticipantService)(nil).Get), ctx, coinProgramParticipantID, userID)
}

// UserList mocks base method.
func (m *MockcoinProgramParticipantService) UserList(ctx context.Context, req dto.CoinProgramParticipantListRequest, userID uuid.UUID) ([]dto.CoinProgramParticipantReturn, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserList", ctx, req, userID)
	ret0, _ := ret[0].([]dto.CoinProgramParticipantReturn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserList indicates an expected call of UserList.
func (mr *MockcoinProgramParticipantServiceMockRecorder) UserList(ctx, req, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserList", reflect.TypeOf((*MockcoinProgramParticipantService)(nil).UserList), ctx, req, userID)
}

// MockrewardService is a mock of rewardService interface.
type MockrewardService struct {
	ctrl     *gomock.Controller
	recorder *MockrewardServiceMockRecorder
	isgomock struct{}
}

// MockrewardServiceMockRecorder is the mock recorder for MockrewardService.
type MockrewardServiceMockRecorder struct {
	mock *MockrewardService
}

// NewMockrewardService creates a new mock instance.
func NewMockrewardService(ctrl *gomock.Controller) *MockrewardService {
	mock := &MockrewardService{ctrl: ctrl}
	mock.recorder = &MockrewardServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockrewardService) EXPECT() *MockrewardServiceMockRecorder {
	return m.recorder
}

// UserListAvailable mocks base method.
func (m *MockrewardService) UserListAvailable(ctx context.Context, req *dto.CoinProgramParticipantListAvailableRequest, userID uuid.UUID) ([]dto.CoinProgramWithRewardsReturn, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserListAvailable", ctx, req, userID)
	ret0, _ := ret[0].([]dto.CoinProgramWithRewardsReturn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserListAvailable indicates an expected call of UserListAvailable.
func (mr *MockrewardServiceMockRecorder) UserListAvailable(ctx, req, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserListAvailable", reflect.TypeOf((*MockrewardService)(nil).UserListAvailable), ctx, req, userID)
}
