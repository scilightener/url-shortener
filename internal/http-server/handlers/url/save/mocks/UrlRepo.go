// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// UrlRepo is an autogenerated mock type for the UrlRepo type
type UrlRepo struct {
	mock.Mock
}

// GetUrl provides a mock function with given fields: alias
func (_m *UrlRepo) GetUrl(alias string) (string, error) {
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

// SaveUrl provides a mock function with given fields: alias, url, validUntil
func (_m *UrlRepo) SaveUrl(alias string, url string, validUntil time.Time) (int64, error) {
	ret := _m.Called(alias, url, validUntil)

	if len(ret) == 0 {
		panic("no return value specified for SaveUrl")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, time.Time) (int64, error)); ok {
		return rf(alias, url, validUntil)
	}
	if rf, ok := ret.Get(0).(func(string, string, time.Time) int64); ok {
		r0 = rf(alias, url, validUntil)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string, string, time.Time) error); ok {
		r1 = rf(alias, url, validUntil)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUrlRepo creates a new instance of UrlRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUrlRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *UrlRepo {
	mock := &UrlRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
