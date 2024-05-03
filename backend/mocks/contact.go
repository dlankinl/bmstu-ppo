// Code generated by MockGen. DO NOT EDIT.
// Source: domain/contact.go
//
// Generated by this command:
//
//	mockgen -source=domain/contact.go -destination=mocks/contact.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	domain "ppo/domain"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockIContactsRepository is a mock of IContactsRepository interface.
type MockIContactsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIContactsRepositoryMockRecorder
}

// MockIContactsRepositoryMockRecorder is the mock recorder for MockIContactsRepository.
type MockIContactsRepositoryMockRecorder struct {
	mock *MockIContactsRepository
}

// NewMockIContactsRepository creates a new mock instance.
func NewMockIContactsRepository(ctrl *gomock.Controller) *MockIContactsRepository {
	mock := &MockIContactsRepository{ctrl: ctrl}
	mock.recorder = &MockIContactsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIContactsRepository) EXPECT() *MockIContactsRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIContactsRepository) Create(ctx context.Context, contact *domain.Contact) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, contact)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockIContactsRepositoryMockRecorder) Create(ctx, contact any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIContactsRepository)(nil).Create), ctx, contact)
}

// DeleteById mocks base method.
func (m *MockIContactsRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockIContactsRepositoryMockRecorder) DeleteById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockIContactsRepository)(nil).DeleteById), ctx, id)
}

// GetById mocks base method.
func (m *MockIContactsRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", ctx, id)
	ret0, _ := ret[0].(*domain.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockIContactsRepositoryMockRecorder) GetById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockIContactsRepository)(nil).GetById), ctx, id)
}

// GetByOwnerId mocks base method.
func (m *MockIContactsRepository) GetByOwnerId(ctx context.Context, id uuid.UUID, page int) ([]*domain.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByOwnerId", ctx, id, page)
	ret0, _ := ret[0].([]*domain.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByOwnerId indicates an expected call of GetByOwnerId.
func (mr *MockIContactsRepositoryMockRecorder) GetByOwnerId(ctx, id, page any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByOwnerId", reflect.TypeOf((*MockIContactsRepository)(nil).GetByOwnerId), ctx, id, page)
}

// Update mocks base method.
func (m *MockIContactsRepository) Update(ctx context.Context, contact *domain.Contact) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, contact)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIContactsRepositoryMockRecorder) Update(ctx, contact any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIContactsRepository)(nil).Update), ctx, contact)
}

// MockIContactsService is a mock of IContactsService interface.
type MockIContactsService struct {
	ctrl     *gomock.Controller
	recorder *MockIContactsServiceMockRecorder
}

// MockIContactsServiceMockRecorder is the mock recorder for MockIContactsService.
type MockIContactsServiceMockRecorder struct {
	mock *MockIContactsService
}

// NewMockIContactsService creates a new mock instance.
func NewMockIContactsService(ctrl *gomock.Controller) *MockIContactsService {
	mock := &MockIContactsService{ctrl: ctrl}
	mock.recorder = &MockIContactsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIContactsService) EXPECT() *MockIContactsServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIContactsService) Create(contact *domain.Contact) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", contact)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockIContactsServiceMockRecorder) Create(contact any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIContactsService)(nil).Create), contact)
}

// DeleteById mocks base method.
func (m *MockIContactsService) DeleteById(id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockIContactsServiceMockRecorder) DeleteById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockIContactsService)(nil).DeleteById), id)
}

// GetById mocks base method.
func (m *MockIContactsService) GetById(id uuid.UUID) (*domain.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", id)
	ret0, _ := ret[0].(*domain.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockIContactsServiceMockRecorder) GetById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockIContactsService)(nil).GetById), id)
}

// GetByOwnerId mocks base method.
func (m *MockIContactsService) GetByOwnerId(id uuid.UUID, page int) ([]*domain.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByOwnerId", id, page)
	ret0, _ := ret[0].([]*domain.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByOwnerId indicates an expected call of GetByOwnerId.
func (mr *MockIContactsServiceMockRecorder) GetByOwnerId(id, page any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByOwnerId", reflect.TypeOf((*MockIContactsService)(nil).GetByOwnerId), id, page)
}

// Update mocks base method.
func (m *MockIContactsService) Update(contact *domain.Contact) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", contact)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIContactsServiceMockRecorder) Update(contact any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIContactsService)(nil).Update), contact)
}
