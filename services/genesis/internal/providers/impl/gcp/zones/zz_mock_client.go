// Code generated by mockery v2.14.1. DO NOT EDIT.

package gcpzones

import (
	mock "github.com/stretchr/testify/mock"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	typedefs "go.taskfleet.io/services/genesis/internal/typedefs"
)

// MockClient is an autogenerated mock type for the Client type
type MockClient struct {
	mock.Mock
}

type MockClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockClient) EXPECT() *MockClient_Expecter {
	return &MockClient_Expecter{mock: &_m.Mock}
}

// GetAccelerator provides a mock function with given fields: zone, kind
func (_m *MockClient) GetAccelerator(zone string, kind typedefs.GPUKind) (Accelerator, error) {
	ret := _m.Called(zone, kind)

	var r0 Accelerator
	if rf, ok := ret.Get(0).(func(string, typedefs.GPUKind) Accelerator); ok {
		r0 = rf(zone, kind)
	} else {
		r0 = ret.Get(0).(Accelerator)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, typedefs.GPUKind) error); ok {
		r1 = rf(zone, kind)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockClient_GetAccelerator_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAccelerator'
type MockClient_GetAccelerator_Call struct {
	*mock.Call
}

// GetAccelerator is a helper method to define mock.On call
//   - zone string
//   - kind typedefs.GPUKind
func (_e *MockClient_Expecter) GetAccelerator(zone interface{}, kind interface{}) *MockClient_GetAccelerator_Call {
	return &MockClient_GetAccelerator_Call{Call: _e.mock.On("GetAccelerator", zone, kind)}
}

func (_c *MockClient_GetAccelerator_Call) Run(run func(zone string, kind typedefs.GPUKind)) *MockClient_GetAccelerator_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(typedefs.GPUKind))
	})
	return _c
}

func (_c *MockClient_GetAccelerator_Call) Return(_a0 Accelerator, _a1 error) *MockClient_GetAccelerator_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetSubnetwork provides a mock function with given fields: zone
func (_m *MockClient) GetSubnetwork(zone string) (string, error) {
	ret := _m.Called(zone)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(zone)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(zone)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockClient_GetSubnetwork_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSubnetwork'
type MockClient_GetSubnetwork_Call struct {
	*mock.Call
}

// GetSubnetwork is a helper method to define mock.On call
//   - zone string
func (_e *MockClient_Expecter) GetSubnetwork(zone interface{}) *MockClient_GetSubnetwork_Call {
	return &MockClient_GetSubnetwork_Call{Call: _e.mock.On("GetSubnetwork", zone)}
}

func (_c *MockClient_GetSubnetwork_Call) Run(run func(zone string)) *MockClient_GetSubnetwork_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockClient_GetSubnetwork_Call) Return(_a0 string, _a1 error) *MockClient_GetSubnetwork_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// List provides a mock function with given fields:
func (_m *MockClient) List() []providers.Zone {
	ret := _m.Called()

	var r0 []providers.Zone
	if rf, ok := ret.Get(0).(func() []providers.Zone); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]providers.Zone)
		}
	}

	return r0
}

// MockClient_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type MockClient_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
func (_e *MockClient_Expecter) List() *MockClient_List_Call {
	return &MockClient_List_Call{Call: _e.mock.On("List")}
}

func (_c *MockClient_List_Call) Run(run func()) *MockClient_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockClient_List_Call) Return(_a0 []providers.Zone) *MockClient_List_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewMockClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockClient creates a new instance of MockClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockClient(t mockConstructorTestingTNewMockClient) *MockClient {
	mock := &MockClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
