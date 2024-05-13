package warehouses_management_system

import (
	"context"

	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_6/orders_management_system/internal/app/models"
	pkgerrors "github.com/Balun-courses/microservices_like_in_bigtech/lecture_6/orders_management_system/pkg/errors"
)

func (r *Client) ReserveStocks(ctx context.Context, userID models.UserID, items []models.Item) error {
	const api = "warehouses_management_system.ReserveStocks"

	/* call exteranl service */

	return pkgerrors.Wrap(api, models.ErrUnimplemented)
	// return nil
}
