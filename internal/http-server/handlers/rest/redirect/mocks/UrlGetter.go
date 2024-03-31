// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UrlGetter is an autogenerated mock type for the UrlGetter type
type UrlGetter struct {
	mock.Mock
}

// GetUrl provides a mock function with given fields: alias
func (_m *UrlGetter) GetUrl(alias string) (string, error) {
	ret := _m.Called(alias)

	if len(ret) == 0 {
		panic("no return value specified for GetUrl")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(alias)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(alias)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(alias)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUrlGetter creates a new instance of UrlGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUrlGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *UrlGetter {
	mock := &UrlGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
