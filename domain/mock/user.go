// Code generated by MockGen. DO NOT EDIT.
// Source: domain/user.go
//
// Generated by this command:
//
//	mockgen -source=domain/user.go -destination=domain/mock/user.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"
	config "togo/config"
	domain "togo/domain"

	gomock "go.uber.org/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// CreateAdmin mocks base method.
func (m *MockUserService) CreateAdmin(cfg *config.Config) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateAdmin", cfg)
}

// CreateAdmin indicates an expected call of CreateAdmin.
func (mr *MockUserServiceMockRecorder) CreateAdmin(cfg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAdmin", reflect.TypeOf((*MockUserService)(nil).CreateAdmin), cfg)
}

// Login mocks base method.
func (m *MockUserService) Login(user domain.User) (string, domain.ResponseError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(domain.ResponseError)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserServiceMockRecorder) Login(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserService)(nil).Login), user)
}

// ParseToken mocks base method.
func (m *MockUserService) ParseToken(token string) (uint, domain.ResponseError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", token)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(domain.ResponseError)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockUserServiceMockRecorder) ParseToken(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockUserService)(nil).ParseToken), token)
}

// SignUp mocks base method.
func (m *MockUserService) SignUp(userSignUp domain.User) domain.ResponseError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", userSignUp)
	ret0, _ := ret[0].(domain.ResponseError)
	return ret0
}

// SignUp indicates an expected call of SignUp.
func (mr *MockUserServiceMockRecorder) SignUp(userSignUp any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockUserService)(nil).SignUp), userSignUp)
}

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserRepository) CreateUser(user domain.User) (domain.User, domain.ResponseError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(domain.ResponseError)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepositoryMockRecorder) CreateUser(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), user)
}

// GetUserByEmail mocks base method.
func (m *MockUserRepository) GetUserByEmail(arg0 string) (domain.User, domain.ResponseError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", arg0)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(domain.ResponseError)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserRepositoryMockRecorder) GetUserByEmail(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUserRepository)(nil).GetUserByEmail), arg0)
}

// GetUserByID mocks base method.
func (m *MockUserRepository) GetUserByID(arg0 uint) (domain.User, domain.ResponseError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", arg0)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(domain.ResponseError)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserRepositoryMockRecorder) GetUserByID(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserRepository)(nil).GetUserByID), arg0)
}

// IsAdmin mocks base method.
func (m *MockUserRepository) IsAdmin(userID uint) (bool, domain.ResponseError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAdmin", userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(domain.ResponseError)
	return ret0, ret1
}

// IsAdmin indicates an expected call of IsAdmin.
func (mr *MockUserRepositoryMockRecorder) IsAdmin(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAdmin", reflect.TypeOf((*MockUserRepository)(nil).IsAdmin), userID)
}
