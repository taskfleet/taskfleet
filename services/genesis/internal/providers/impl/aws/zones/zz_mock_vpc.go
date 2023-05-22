// Code generated by mockery v2.13.1. DO NOT EDIT.

package awszones

import (
	context "context"

	ec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	mock "github.com/stretchr/testify/mock"
)

// MockVPC is an autogenerated mock type for the VPC type
type MockVPC struct {
	mock.Mock
}

// Client provides a mock function with given fields:
func (_m *MockVPC) Client() *ec2.Client {
	ret := _m.Called()

	var r0 *ec2.Client
	if rf, ok := ret.Get(0).(func() *ec2.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ec2.Client)
		}
	}

	return r0
}

// ID provides a mock function with given fields:
func (_m *MockVPC) ID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Region provides a mock function with given fields:
func (_m *MockVPC) Region() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SecurityGroups provides a mock function with given fields: ctx
func (_m *MockVPC) SecurityGroups(ctx context.Context) ([]string, error) {
	ret := _m.Called(ctx)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context) []string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Subnets provides a mock function with given fields: ctx
func (_m *MockVPC) Subnets(ctx context.Context) (map[string]string, error) {
	ret := _m.Called(ctx)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(context.Context) map[string]string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockVPC interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockVPC creates a new instance of MockVPC. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockVPC(t mockConstructorTestingTNewMockVPC) *MockVPC {
	mock := &MockVPC{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}