package orders_management_system

import (
	"context"
	"errors"

	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_6/orders_management_system/internal/app/models"
	pkgerrors "github.com/Balun-courses/microservices_like_in_bigtech/lecture_6/orders_management_system/pkg/errors"
	"github.com/google/uuid"
)

// CreateOrder - создание заказа
func (oms *usecase) CreateOrder(ctx context.Context, userID models.UserID, info CreateOrderInfo) (*models.Order, error) {
	const api = "orders_management_system.usecase.CreateOrder"

	// Резервируем стоки на складах
	if err := oms.WarehouseManagementSystem.ReserveStocks(ctx, userID, info.Items); err != nil {
		return nil, pkgerrors.Wrap(api, err)
	}

	//

	// Формируем запись о заказе
	var (
		orderID = models.OrderID(uuid.New())
		order   = &models.Order{
			ID:                orderID,
			UserID:            userID,
			Items:             info.Items,
			DeliveryOrderInfo: info.DeliveryOrderInfo,
		}
	)

	const retries = 3
	var err error
	for i := 1; i <= retries; i++ { // ретраи
		// Создаем заказ в БД (Postgres)
		if err = oms.WriteOrdersStorage.CreateOrder(ctx, order); err != nil {
			if errors.Is(err, models.ErrAlreadyExists) {
				order.ID = models.OrderID(uuid.New())
			}
			continue
		}
		break
	}
	if err != nil {
		return nil, pkgerrors.Wrap(api, err)
	}

	go oms.OrderAnalyticsUsecaseInterface.SendOrderData(ctx, order)

	/*
		-> [OMS] -> [DB]
		     |
			  -> [Kafka] -> [Analytic-consumer] -> [Clickhouse] -> отчет
	*/
	/*
		-> [OMS] -> [KV] -> [Postgres] -> jn

	*/

	return order, nil
}
