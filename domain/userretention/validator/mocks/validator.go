// Code generated by mockery v2.6.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UserRetentionValidator is an autogenerated mock type for the UserRetentionValidator type
type UserRetentionValidator struct {
	mock.Mock
}

// ValidateInput provides a mock function with given fields: filePath
func (_m *UserRetentionValidator) ValidateInput(filePath string) error {
	ret := _m.Called(filePath)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(filePath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
