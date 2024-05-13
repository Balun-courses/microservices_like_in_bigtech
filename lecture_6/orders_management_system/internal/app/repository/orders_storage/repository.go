package orders_storage

import (
	oms "github.com/Balun-courses/microservices_like_in_bigtech/lecture_6/orders_management_system/internal/app/usecases/orders_management_system"
)

type OrdersStorage struct {
	/*
		PostgreSQL, MSSQL, MySQL, Redis, any you want...
	*/
}

// Check that we implemet contract for usecase
var (
	// oms.OrdersStorage - PORT
	_ oms.ReadOrdersStorage  = (*OrdersStorage)(nil)
	_ oms.WriteOrdersStorage = (*OrdersStorage)(nil)
)

// New - returns OrdersStorage
func New( /* ... */ ) *OrdersStorage {
	return &OrdersStorage{
		/* ... */
	}
}
