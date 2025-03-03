// Code generated by mockery v2.33.2. DO NOT EDIT.

package coremocks

import (
	core "github.com/hyperledger/firefly/pkg/core"
	mock "github.com/stretchr/testify/mock"
)

// OperationCallbacks is an autogenerated mock type for the OperationCallbacks type
type OperationCallbacks struct {
	mock.Mock
}

// OperationUpdate provides a mock function with given fields: update
func (_m *OperationCallbacks) OperationUpdate(update *core.OperationUpdate) {
	_m.Called(update)
}

// NewOperationCallbacks creates a new instance of OperationCallbacks. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOperationCallbacks(t interface {
	mock.TestingT
	Cleanup(func())
}) *OperationCallbacks {
	mock := &OperationCallbacks{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
