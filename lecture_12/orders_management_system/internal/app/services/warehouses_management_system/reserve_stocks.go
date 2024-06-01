package warehouses_management_system

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/moguchev/microservices_courcse/orders_management_system/internal/app/models"
	pkgerrors "github.com/moguchev/microservices_courcse/orders_management_system/pkg/errors"
	"github.com/moguchev/microservices_courcse/orders_management_system/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

func (r *Client) ReserveStocks(
	ctx context.Context,
	userID models.UserID,
	items []models.Item,
) error {
	const api = "warehouses_management_system.ReserveStocks"

	span, ctx := opentracing.StartSpanFromContext(ctx, "warehouses_management_system.ReserveStocks")
	defer span.Finish()

	span.SetTag("user_id", userID)

	logger.Info(ctx, "stock reserved")

	/* call external service */
	time.Sleep(50 * time.Millisecond)

	if rand.Int31n(4) == 3 {
		return pkgerrors.Wrap(api, errors.New("error"))
	}
	return nil
	// return pkgerrors.Wrap(api, models.ErrUnimplemented)
}
