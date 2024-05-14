package models

import (
	"time"
)

// Order - заказ
type Order struct {
	ID                OrderID // ID заказа
	UserID            UserID  // ID пользователя (чей заказ)
	Items             []Item  // Информация о составе заказа
	DeliveryOrderInfo         // Информация о доставке
	/* ... */
}

// DeliveryOrderInfo - информация о доставке заказа
type DeliveryOrderInfo struct {
	DeliveryVariantID DeliveryVariantID
	DeliveryDate      time.Time
}

// Item - информация о составе заказа
type Item struct {
	SKU         SKU         // SKU
	Quantity    uint32      // количество SKU
	WarehouseID WarehouseID // с какого склада будет браться сток
}
