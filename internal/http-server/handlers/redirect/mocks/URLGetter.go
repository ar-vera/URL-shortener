package mocks

import mock "github.com/stretchr/testify/mock"

type URLGetter struct {
	mock.Mock
}

func (_m *URLGetter) GetURL(alias string) (string, error) {
	ret := _m.Called(alias)

	var r0 string
	var r1 error

	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(alias)
	}

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(string)
	}

	r1 = ret.Error(1)

	return r0, r1
}

type mockConstructorTestingTNewURLGetter interface {
	mock.TestingT
	Cleanup(func())
}

func NewURLGetter(t mockConstructorTestingTNewURLGetter) *URLGetter {
	mock := &URLGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
