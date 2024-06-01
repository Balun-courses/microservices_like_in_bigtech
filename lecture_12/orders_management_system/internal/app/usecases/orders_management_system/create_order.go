package orders_management_system

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/moguchev/microservices_courcse/orders_management_system/internal/app/models"
	pkgerrors "github.com/moguchev/microservices_courcse/orders_management_system/pkg/errors"
	"github.com/moguchev/microservices_courcse/orders_management_system/pkg/transaction_manager"
)

// CreateOrder - создание заказа
func (oms *usecase) CreateOrder(ctx context.Context, userID models.UserID, info CreateOrderInfo) (*models.Order, error) {
	const api = "orders_management_system.usecase.CreateOrder"

	// TODO: ключ идемпотентности

	// Резервируем стоки на складах
	if err := oms.WarehouseManagementSystem.ReserveStocks(ctx, userID, info.Items); err != nil {
		return nil, pkgerrors.Wrap(api, err)
	}

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
		err := oms.TransactionManager.RunReadCommitted(ctx, transaction_manager.ReadWrite,
			func(txCtx context.Context) error { // TRANSANCTION SCOPE
				// Создаем заказ в БД
				if err = oms.OrdersStorage.CreateOrder(txCtx, order); err != nil {
					return err
				}
				// Создаем сообщение outbox в БД
				if err = oms.CreateOutboxMessage(txCtx, order); err != nil {
					return err
				}

				return nil
			},
		)
		if err != nil {
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

	return order, nil
}
