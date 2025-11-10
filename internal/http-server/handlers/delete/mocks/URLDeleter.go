package mocks

import (
	"github.com/stretchr/testify/mock"
)

type URLDeleter struct {
	mock.Mock
}

func (_m *URLDeleter) DeleteURL(alias string) error {
	ret := _m.Called(alias)

	var r0 error

	if rf, ok := ret.Get(0).(func(string) error); ok {
		return rf(alias)
	}

	r0 = ret.Error(0)

	return r0
}

type mockConstructorTestingTNewURLDeleter interface {
	mock.TestingT
	Cleanup(func())
}

func NewURLDeleter(t mockConstructorTestingTNewURLDeleter) *URLDeleter {
	mock := &URLDeleter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
