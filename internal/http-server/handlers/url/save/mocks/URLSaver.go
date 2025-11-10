package mocks

import (
	"github.com/stretchr/testify/mock"
)

type URLSaver struct {
	mock.Mock
}

func (_m *URLSaver) SaveURL(urlToSave string, alias string) (int64, error) {
	ret := _m.Called(urlToSave, alias)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (int64, error)); ok {
		r0, r1 = rf(urlToSave, alias)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(int64)
		}
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewURLSaver interface {
	mock.TestingT
	Cleanup(func())
}

func NewURLSaver(t mockConstructorTestingTNewURLSaver) *URLSaver {
	mock := &URLSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
