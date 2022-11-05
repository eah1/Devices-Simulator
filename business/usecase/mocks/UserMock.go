// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	webmodels "device-simulator/business/web/webmodels"
)

// UserMock is an autogenerated mock type for the User type
type UserMock struct {
	mock.Mock
}

// RegisterUser provides a mock function with given fields: userRegister
func (_m *UserMock) RegisterUser(userRegister webmodels.RegisterUser) error {
	ret := _m.Called(userRegister)

	var r0 error
	if rf, ok := ret.Get(0).(func(webmodels.RegisterUser) error); ok {
		r0 = rf(userRegister)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUserMock interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserMock creates a new instance of UserMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserMock(t mockConstructorTestingTNewUserMock) *UserMock {
	mock := &UserMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
