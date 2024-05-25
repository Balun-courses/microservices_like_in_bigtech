package usecases

import (
	"context"

	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_7/orders_management_system/internal/app/models"
)

// OMSUsecaseInterface - интерфейс бизнес логики
type OMSUsecaseInterface interface {
	// CreateOrder - создание заказа
	//
	// @errors: ErrReserveStocks
	CreateOrder(ctx context.Context, userID models.UserID, info CreateOrderInfo) (*models.Order, error)
}

// VSCODE
// Goland

// gRPC
// protoc
// buf.buld - 403

// echo, chi, gorrilamx

// POST
// DataGrip /  DBera

// ktx kns

// helm
// ansible/terra
