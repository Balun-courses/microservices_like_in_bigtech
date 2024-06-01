package orders_storage

import (
	"database/sql"
	"encoding/json"

	"github.com/moguchev/microservices_courcse/orders_management_system/internal/app/models"
	pkgerrors "github.com/moguchev/microservices_courcse/orders_management_system/pkg/errors"
	uuid "github.com/vgarvardt/pgx-google-uuid/v5"
)

// модель представления данных в БД (храним в json)
type orderItem struct {
	SKUID       int64 `json:"sku_id"`       // id товара
	Quantity    int32 `json:"quantity"`     // количество SKU
	WarehouseID int64 `json:"warehouse_id"` // с какого склада будет браться сток
}

func getOrderItems(order *models.Order) []orderItem {
	items := make([]orderItem, len(order.Items))
	for i := range order.Items {
		items[i] = orderItem{
			SKUID:       int64(order.Items[i].SKU.ID),
			Quantity:    int32(order.Items[i].Quantity),
			WarehouseID: int64(order.Items[i].WarehouseID),
		}
	}
	return items
}

type orderRow struct {
	ID                uuid.UUID     `db:"id"`
	UserID            int64         `db:"user_id"`
	Items             []byte        `db:"items"`
	DeliveryVariantID sql.NullInt64 `db:"delivery_variant_id"`
	DeliveryDate      sql.NullTime  `db:"delivery_date"`
}

func (r *orderRow) ValuesMap() map[string]any {
	return map[string]any{
		"id":                  r.ID,
		"user_id":             r.UserID,
		"items":               r.Items,
		"delivery_variant_id": r.DeliveryVariantID,
		"delivery_date":       r.DeliveryDate,
	}
}

func (r *orderRow) Values(columns ...string) []any {
	values := make([]any, 0, len(columns))
	m := r.ValuesMap()

	for i := range columns {
		values = append(values, m[columns[i]])
	}

	return values
}

func newOrderRowFromModelsOrder(order *models.Order) (*orderRow, error) {
	items, err := json.Marshal(getOrderItems(order))
	if err != nil {
		return nil, pkgerrors.Wrap("newOrderRowFromModelsOrder", err)
	}

	return &orderRow{
		ID:     uuid.UUID(order.ID),
		UserID: int64(order.UserID),
		Items:  items,
		DeliveryVariantID: sql.NullInt64{
			Int64: int64(order.DeliveryVariantID),
			Valid: order.DeliveryVariantID != 0,
		},
		DeliveryDate: sql.NullTime{
			Time:  order.DeliveryDate,
			Valid: order.DeliveryDate.IsZero(),
		},
	}, nil
}
