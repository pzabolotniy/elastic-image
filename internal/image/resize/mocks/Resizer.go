// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import io "io"
import logging "github.com/pzabolotniy/elastic-image/internal/logging"
import mock "github.com/stretchr/testify/mock"

// Resizer is an autogenerated mock type for the Resizer type
type Resizer struct {
	mock.Mock
}

// Logger provides a mock function with given fields:
func (_m *Resizer) Logger() logging.Logger {
	ret := _m.Called()

	var r0 logging.Logger
	if rf, ok := ret.Get(0).(func() logging.Logger); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(logging.Logger)
		}
	}

	return r0
}

// Resize provides a mock function with given fields: src, width, height
func (_m *Resizer) Resize(src io.Reader, width uint, height uint) ([]byte, error) {
	ret := _m.Called(src, width, height)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(io.Reader, uint, uint) []byte); ok {
		r0 = rf(src, width, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(io.Reader, uint, uint) error); ok {
		r1 = rf(src, width, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}