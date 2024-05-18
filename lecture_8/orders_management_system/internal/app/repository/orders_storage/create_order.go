package orders_storage

import (
	"context"
	"errors"

	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_8/orders_management_system/internal/app/models"
	pkgerrors "github.com/Balun-courses/microservices_like_in_bigtech/lecture_8/orders_management_system/pkg/errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *OrdersStorage) CreateOrder(ctx context.Context, order *models.Order) error {
	const api = "orders_storage.CreateOrder"

	row, err := newOrderRowFromModelsOrder(order)
	if err != nil {
		return pkgerrors.Wrap(api, err)
	}

	columns := []string{
		"id",                  // uuid
		"user_id",             // int8
		"items",               // json
		"delivery_variant_id", // int8
		"delivery_date",       // int8
	}

	// вариант 1
	query := squirrel.Insert(tableOrdersName).
		Columns(columns...).
		Values(row.Values(columns...)...).
		PlaceholderFormat(squirrel.Dollar)

	// вариант 2
	query = squirrel.Insert(tableOrdersName).
		SetMap(row.ValuesMap()).
		PlaceholderFormat(squirrel.Dollar)

	if _, err := r.driver.GetQueryEngine(ctx).Execx(ctx, query); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
			return pkgerrors.Wrap(api, models.ErrAlreadyExists)
		}
		return pkgerrors.Wrap(api, err)
	}

	return nil
}
