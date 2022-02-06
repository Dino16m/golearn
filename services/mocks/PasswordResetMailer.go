// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// PasswordResetMailer is an autogenerated mock type for the PasswordResetMailer type
type PasswordResetMailer struct {
	mock.Mock
}

// SendResetCode provides a mock function with given fields: name, email, code
func (_m *PasswordResetMailer) SendResetCode(name string, email string, code int) error {
	ret := _m.Called(name, email, code)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, int) error); ok {
		r0 = rf(name, email, code)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendResetLink provides a mock function with given fields: name, email, link
func (_m *PasswordResetMailer) SendResetLink(name string, email string, link string) error {
	ret := _m.Called(name, email, link)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(name, email, link)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
