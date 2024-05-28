// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/Balun-courses/microservices_like_in_bigtech/lecture_11/orders_management_system/internal/app/models"
	mock "github.com/stretchr/testify/mock"
)

// OrdersStorage is an autogenerated mock type for the OrdersStorage type
type OrdersStorage struct {
	mock.Mock
}

// CreateOrder provides a mock function with given fields: ctx, order
func (_m *OrdersStorage) CreateOrder(ctx context.Context, order *models.Order) error {
	ret := _m.Called(ctx, order)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Order) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewOrdersStorage creates a new instance of OrdersStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrdersStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrdersStorage {
	mock := &OrdersStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
