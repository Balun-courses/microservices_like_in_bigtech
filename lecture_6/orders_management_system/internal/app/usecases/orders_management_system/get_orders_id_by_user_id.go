package orders_management_system

import (
	"context"

	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_6/orders_management_system/internal/app/models"
)

func (oms *usecase) GetOrdersIDByUserID(ctx context.Context, userID models.UserID, status *models.OrderStatus) ([]models.OrderID, error) {
	return nil, models.ErrUnimplemented
}
