package orders_storage

import (
	"context"

	"github.com/moguchev/microservices_courcse/orders_management_system/internal/app/models"
	uuid "github.com/vgarvardt/pgx-google-uuid/v5"
)

func (r *OrdersStorage) CreateOutboxMessage(ctx context.Context, order *models.Order) error {
	if _, err := r.driver.GetQueryEngine(ctx).Exec(ctx,
		"INSERT INTO orders_outbox_messages (order_id) VALUES ($1)",
		uuid.UUID(order.ID),
	); err != nil {
		return err
	}
	return nil
}
