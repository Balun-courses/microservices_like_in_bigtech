package orders_storage

import (
	"context"

	oms "github.com/Balun-courses/microservices_like_in_bigtech/lecture_11/orders_management_system/internal/app/usecases/orders_management_system"
	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_11/orders_management_system/pkg/postgres"
	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_11/orders_management_system/pkg/transaction_manager"
	"github.com/jackc/pgx/v5/pgconn"
)

// Check that we implemet contract for usecase
var (
	_ oms.OrdersStorage = (*OrdersStorage)(nil)
)

// для удобства тестирования (мока базы)
type Connection interface {
	Execx(ctx context.Context, sqlizer postgres.Sqlizer) (pgconn.CommandTag, error)
}

type OrdersStorage struct {
	// connection *postgres.Connection // если тетсируте только интеграционными
	// connection Connection // если мокаете базу данных
	driver QueryEngineProvider
}

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) transaction_manager.QueryEngine
}

// New - returns OrdersStorage
func New( /*connection *postgres.Connection*/ driver QueryEngineProvider) *OrdersStorage {
	return &OrdersStorage{
		// connection: connection,
		driver: driver,
	}
}

const (
	tableOrdersName = "orders"
)
