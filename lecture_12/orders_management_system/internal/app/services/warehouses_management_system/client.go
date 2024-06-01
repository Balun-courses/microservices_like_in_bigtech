package warehouses_management_system

import "github.com/moguchev/microservices_courcse/orders_management_system/internal/app/usecases/orders_management_system"

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
