package orders_management_system

import "github.com/moguchev/microservices_courcse/orders_management_system/internal/app/models"

// CreateOrderInputInfo - DTO заказа (для создания заказа)
type CreateOrderInfo struct {
	Items             []models.Item            // Товары в заказе
	DeliveryOrderInfo models.DeliveryOrderInfo // Информация о доставке
}
