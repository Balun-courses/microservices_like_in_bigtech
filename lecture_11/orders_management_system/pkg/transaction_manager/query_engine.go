package transaction_manager

import (
	"context"

	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_11/orders_management_system/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// QueryEngineProvider - smths that gives us QueryEngine
type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) QueryEngine
}

// PgxCommonAPI - pgx common api
type PgxCommonAPI interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// PgxCommonScanAPI улучшенный PgxCommonAPI
type PgxCommonScanAPI interface {
	// Getx - aka QueryRow
	Getx(ctx context.Context, dest interface{}, sqlizer postgres.Sqlizer) error
	// Selectx - aka Query
	Selectx(ctx context.Context, dest interface{}, sqlizer postgres.Sqlizer) error
	// Execx - aka Exec
	Execx(ctx context.Context, sqlizer postgres.Sqlizer) (pgconn.CommandTag, error)
}

// PgxExtendedAPI - ...
type PgxExtendedAPI interface {
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}

// QueryEngine is a common database query interface.
type QueryEngine interface {
	PgxCommonAPI
	PgxCommonScanAPI
	PgxExtendedAPI
}

// TxAccessMode is the transaction access mode (read write or read only)
type TxAccessMode = pgx.TxAccessMode

// Transaction access modes
const (
	ReadWrite = pgx.ReadWrite
	ReadOnly  = pgx.ReadOnly
)
