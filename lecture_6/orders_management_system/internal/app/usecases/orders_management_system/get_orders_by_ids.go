package orders_management_system

import (
	"context"

	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_6/orders_management_system/internal/app/models"
)

func (oms *usecase) GetOrdersByIDs(ctx context.Context, ids []models.OrderID) ([]*models.Order, error) {
	return nil, models.ErrUnimplemented
}
