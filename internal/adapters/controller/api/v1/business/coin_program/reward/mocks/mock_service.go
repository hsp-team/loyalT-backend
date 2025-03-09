// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go
//
// Generated by this command:
//
//	mockgen -source=handler.go -destination=mocks/mock_service.go -package=mocks rewardService
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

// Create mocks base method.
func (m *MockrewardService) Create(ctx context.Context, req *dto.RewardCreateRequest, businessID uuid.UUID) (*dto.RewardCreateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, req, businessID)
	ret0, _ := ret[0].(*dto.RewardCreateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockrewardServiceMockRecorder) Create(ctx, req, businessID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockrewardService)(nil).Create), ctx, req, businessID)
}

// Delete mocks base method.
func (m *MockrewardService) Delete(ctx context.Context, req *dto.RewardDeleteRequest, businessID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, req, businessID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockrewardServiceMockRecorder) Delete(ctx, req, businessID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockrewardService)(nil).Delete), ctx, req, businessID)
}

// List mocks base method.
func (m *MockrewardService) List(ctx context.Context, request *dto.RewardListRequest, businessID uuid.UUID) ([]dto.RewardReturn, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, request, businessID)
	ret0, _ := ret[0].([]dto.RewardReturn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockrewardServiceMockRecorder) List(ctx, request, businessID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockrewardService)(nil).List), ctx, request, businessID)
}

// MockqrService is a mock of qrService interface.
type MockqrService struct {
	ctrl     *gomock.Controller
	recorder *MockqrServiceMockRecorder
	isgomock struct{}
}

// MockqrServiceMockRecorder is the mock recorder for MockqrService.
type MockqrServiceMockRecorder struct {
	mock *MockqrService
}

// NewMockqrService creates a new mock instance.
func NewMockqrService(ctrl *gomock.Controller) *MockqrService {
	mock := &MockqrService{ctrl: ctrl}
	mock.recorder = &MockqrServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockqrService) EXPECT() *MockqrServiceMockRecorder {
	return m.recorder
}

// ActivateUserReward mocks base method.
func (m *MockqrService) ActivateUserReward(ctx context.Context, req *dto.RewardActivateRequest, businessID uuid.UUID) (*dto.RewardReturn, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActivateUserReward", ctx, req, businessID)
	ret0, _ := ret[0].(*dto.RewardReturn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ActivateUserReward indicates an expected call of ActivateUserReward.
func (mr *MockqrServiceMockRecorder) ActivateUserReward(ctx, req, businessID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActivateUserReward", reflect.TypeOf((*MockqrService)(nil).ActivateUserReward), ctx, req, businessID)
}
