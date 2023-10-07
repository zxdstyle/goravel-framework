// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	queue "github.com/goravel/framework/contracts/queue"
	mock "github.com/stretchr/testify/mock"
)

// Driver is an autogenerated mock type for the Driver type
type Driver struct {
	mock.Mock
}

// Bulk provides a mock function with given fields: jobs
func (_m *Driver) Bulk(jobs []queue.Job) error {
	ret := _m.Called(jobs)

	var r0 error
	if rf, ok := ret.Get(0).(func([]queue.Job) error); ok {
		r0 = rf(jobs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Clear provides a mock function with given fields:
func (_m *Driver) Clear() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ConnectionName provides a mock function with given fields:
func (_m *Driver) ConnectionName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Delete provides a mock function with given fields: job
func (_m *Driver) Delete(job queue.Job) error {
	ret := _m.Called(job)

	var r0 error
	if rf, ok := ret.Get(0).(func(queue.Job) error); ok {
		r0 = rf(job)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Later provides a mock function with given fields: job, delay
func (_m *Driver) Later(job queue.Job, delay int) error {
	ret := _m.Called(job, delay)

	var r0 error
	if rf, ok := ret.Get(0).(func(queue.Job, int) error); ok {
		r0 = rf(job, delay)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Pop provides a mock function with given fields:
func (_m *Driver) Pop() (queue.Job, error) {
	ret := _m.Called()

	var r0 queue.Job
	var r1 error
	if rf, ok := ret.Get(0).(func() (queue.Job, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() queue.Job); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(queue.Job)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Push provides a mock function with given fields: job
func (_m *Driver) Push(job queue.Job) error {
	ret := _m.Called(job)

	var r0 error
	if rf, ok := ret.Get(0).(func(queue.Job) error); ok {
		r0 = rf(job)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Release provides a mock function with given fields: job, delay
func (_m *Driver) Release(job queue.Job, delay int) error {
	ret := _m.Called(job, delay)

	var r0 error
	if rf, ok := ret.Get(0).(func(queue.Job, int) error); ok {
		r0 = rf(job, delay)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Size provides a mock function with given fields:
func (_m *Driver) Size() (int, error) {
	ret := _m.Called()

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func() (int, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewDriver creates a new instance of Driver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDriver(t interface {
	mock.TestingT
	Cleanup(func())
}) *Driver {
	mock := &Driver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
