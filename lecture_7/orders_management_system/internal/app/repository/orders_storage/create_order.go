package orders_storage

import (
	"context"

	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_7/orders_management_system/internal/app/models"
	pkgerrors "github.com/Balun-courses/microservices_like_in_bigtech/lecture_7/orders_management_system/pkg/errors"
)

func (r *OrdersStorage) CreateOrder(ctx context.Context, order *models.Order) error {
	const api = "orders_storage.CreateOrder"

	/* here your queries */

	return pkgerrors.Wrap(api, models.ErrUnimplemented)
	// return pkgerrors.Wrap(api, models.ErrAlreadyExists)
	// return nil
}
