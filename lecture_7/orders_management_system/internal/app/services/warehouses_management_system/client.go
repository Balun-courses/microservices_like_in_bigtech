package warehouses_management_system

import "github.com/Balun-courses/microservices_like_in_bigtech/lecture_7/orders_management_system/internal/app/usecases/orders_management_system"

type Client struct {
	/*
		HTTP, gRPC, ... client
	*/
}

// Check that we implemet contract for usecase
var _ orders_management_system.WarehouseManagementSystem = (*Client)(nil)

// NewClient - returns WMS service adapter
func NewClient( /* ... */ ) *Client {
	return &Client{
		/* ... */
	}
}
