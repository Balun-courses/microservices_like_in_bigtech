// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/moguchev/microservices_courcse/orders_management_system/internal/app/models"
	mock "github.com/stretchr/testify/mock"
)

// WarehouseManagementSystem is an autogenerated mock type for the WarehouseManagementSystem type
type WarehouseManagementSystem struct {
	mock.Mock
}

// ReserveStocks provides a mock function with given fields: ctx, userID, items
func (_m *WarehouseManagementSystem) ReserveStocks(ctx context.Context, userID models.UserID, items []models.Item) error {
	ret := _m.Called(ctx, userID, items)

	if len(ret) == 0 {
		panic("no return value specified for ReserveStocks")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID, []models.Item) error); ok {
		r0 = rf(ctx, userID, items)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewWarehouseManagementSystem creates a new instance of WarehouseManagementSystem. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWarehouseManagementSystem(t interface {
	mock.TestingT
	Cleanup(func())
}) *WarehouseManagementSystem {
	mock := &WarehouseManagementSystem{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
