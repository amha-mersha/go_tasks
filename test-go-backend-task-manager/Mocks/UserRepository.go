// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/domains"
	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: cxt, newUser
func (_m *UserRepository) CreateUser(cxt context.Context, newUser domain.User) (string, *domain.UserError) {
	ret := _m.Called(cxt, newUser)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 string
	var r1 *domain.UserError
	if rf, ok := ret.Get(0).(func(context.Context, domain.User) (string, *domain.UserError)); ok {
		return rf(cxt, newUser)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.User) string); ok {
		r0 = rf(cxt, newUser)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.User) *domain.UserError); ok {
		r1 = rf(cxt, newUser)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*domain.UserError)
		}
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: cxt, userID
func (_m *UserRepository) DeleteUser(cxt context.Context, userID string) (domain.User, *domain.UserError) {
	ret := _m.Called(cxt, userID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUser")
	}

	var r0 domain.User
	var r1 *domain.UserError
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.User, *domain.UserError)); ok {
		return rf(cxt, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.User); ok {
		r0 = rf(cxt, userID)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) *domain.UserError); ok {
		r1 = rf(cxt, userID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*domain.UserError)
		}
	}

	return r0, r1
}

// FetchAllUsers provides a mock function with given fields: cxt
func (_m *UserRepository) FetchAllUsers(cxt context.Context) ([]domain.User, *domain.UserError) {
	ret := _m.Called(cxt)

	if len(ret) == 0 {
		panic("no return value specified for FetchAllUsers")
	}

	var r0 []domain.User
	var r1 *domain.UserError
	if rf, ok := ret.Get(0).(func(context.Context) ([]domain.User, *domain.UserError)); ok {
		return rf(cxt)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []domain.User); ok {
		r0 = rf(cxt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) *domain.UserError); ok {
		r1 = rf(cxt)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*domain.UserError)
		}
	}

	return r0, r1
}

// FetchUserByID provides a mock function with given fields: cxt, ID
func (_m *UserRepository) FetchUserByID(cxt context.Context, ID string) (domain.User, *domain.UserError) {
	ret := _m.Called(cxt, ID)

	if len(ret) == 0 {
		panic("no return value specified for FetchUserByID")
	}

	var r0 domain.User
	var r1 *domain.UserError
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.User, *domain.UserError)); ok {
		return rf(cxt, ID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.User); ok {
		r0 = rf(cxt, ID)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) *domain.UserError); ok {
		r1 = rf(cxt, ID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*domain.UserError)
		}
	}

	return r0, r1
}

// FetchUserByUsername provides a mock function with given fields: cxt, username
func (_m *UserRepository) FetchUserByUsername(cxt context.Context, username string) (domain.User, *domain.UserError) {
	ret := _m.Called(cxt, username)

	if len(ret) == 0 {
		panic("no return value specified for FetchUserByUsername")
	}

	var r0 domain.User
	var r1 *domain.UserError
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.User, *domain.UserError)); ok {
		return rf(cxt, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.User); ok {
		r0 = rf(cxt, username)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) *domain.UserError); ok {
		r1 = rf(cxt, username)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*domain.UserError)
		}
	}

	return r0, r1
}

// FetchUserCount provides a mock function with given fields: cxt
func (_m *UserRepository) FetchUserCount(cxt context.Context) (int, *domain.UserError) {
	ret := _m.Called(cxt)

	if len(ret) == 0 {
		panic("no return value specified for FetchUserCount")
	}

	var r0 int
	var r1 *domain.UserError
	if rf, ok := ret.Get(0).(func(context.Context) (int, *domain.UserError)); ok {
		return rf(cxt)
	}
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(cxt)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context) *domain.UserError); ok {
		r1 = rf(cxt)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*domain.UserError)
		}
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: cxt, updateUser
func (_m *UserRepository) UpdateUser(cxt context.Context, updateUser domain.User) (domain.User, *domain.UserError) {
	ret := _m.Called(cxt, updateUser)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 domain.User
	var r1 *domain.UserError
	if rf, ok := ret.Get(0).(func(context.Context, domain.User) (domain.User, *domain.UserError)); ok {
		return rf(cxt, updateUser)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.User) domain.User); ok {
		r0 = rf(cxt, updateUser)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.User) *domain.UserError); ok {
		r1 = rf(cxt, updateUser)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*domain.UserError)
		}
	}

	return r0, r1
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
