package models

import "github.com/google/uuid"

// OrderID UUID заказа
type OrderID uuid.UUID

// String - represent OrderID as string
func (v OrderID) String() string {
	return uuid.UUID(v).String()
}

// UserID - тип id пользователя
type UserID uint64

// SKUID - тип id товарной еденицы
type SKUID uint64

// WarehouseID - тип id склада
type WarehouseID uint64

// DeliveryVariantID - тип id способа доставки
type DeliveryVariantID uint64
