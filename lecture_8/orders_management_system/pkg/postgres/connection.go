package postgres

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

const (
	maxConnIdleTimeDefault     = time.Minute
	maxConnLifeTimeDefault     = time.Hour
	minConnectionsCountDefault = 2
	maxConnectionsCountDefault = 10
)

type connectionPoolOptions struct {
	maxConnIdleTime     time.Duration
	maxConnLifeTime     time.Duration
	minConnectionsCount int32
	maxConnectionsCount int32
	tlsConfig           *tls.Config
}

type ConnectionPoolOption func(options *connectionPoolOptions)

// WithMaxConnIdleTime ...
func WithMaxConnIdleTime(d time.Duration) ConnectionPoolOption {
	return func(opts *connectionPoolOptions) {
		opts.maxConnIdleTime = d
	}
}

// WithMaxConnLifeTime ...
func WithMaxConnLifeTime(d time.Duration) ConnectionPoolOption {
	return func(opts *connectionPoolOptions) {
		opts.maxConnLifeTime = d
	}
}

// WithMinConnectionsCount ...
func WithMinConnectionsCount(c int32) ConnectionPoolOption {
	return func(opts *connectionPoolOptions) {
		opts.minConnectionsCount = c
	}
}

// WithMaxConnectionsCount ...
func WithMaxConnectionsCount(c int32) ConnectionPoolOption {
	return func(opts *connectionPoolOptions) {
		opts.maxConnectionsCount = c
	}
}

// WithSSL ...
func WithSSL(cfg *tls.Config) ConnectionPoolOption {
	return func(opts *connectionPoolOptions) {
		opts.tlsConfig = cfg
	}
}

// Connection - postgres connection pool
type Connection struct {
	pool *pgxpool.Pool
}

// NewConnectionPool - returns new Connection (connection pool for postgres)
func NewConnectionPool(ctx context.Context, connString string, opts ...ConnectionPoolOption) (*Connection, error) {
	// parse connString
	connConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("can't parse connection string to config: %w", err)
	}

	connConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	// make options
	options := &connectionPoolOptions{
		maxConnIdleTime:     maxConnIdleTimeDefault,
		maxConnLifeTime:     maxConnLifeTimeDefault,
		minConnectionsCount: minConnectionsCountDefault,
		maxConnectionsCount: maxConnectionsCountDefault,
	}
	for _, opt := range opts {
		opt(options)
	}

	// apply options
	connConfig.MaxConnIdleTime = options.maxConnIdleTime
	connConfig.MaxConnLifetime = options.maxConnLifeTime
	connConfig.MinConns = options.minConnectionsCount
	connConfig.MaxConns = options.maxConnectionsCount
	connConfig.ConnConfig.Config.TLSConfig = options.tlsConfig

	// connect to database
	p, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	// ping database
	if err := p.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database error: %w", err)
	}

	return &Connection{
		pool: p,
	}, nil
}

func (c *Connection) Close() error {
	c.pool.Close()
	return nil
}

// Query - pgx.Query
func (c *Connection) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return c.pool.Query(ctx, sql, args...)
}

// Query - pgx.Exec
func (c *Connection) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return c.pool.Exec(ctx, sql, args...)
}

// Query - pgx.QueryRow
func (c *Connection) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return c.pool.QueryRow(ctx, sql, args...)
}

// Begin - pgx.Begin
func (c *Connection) Begin(ctx context.Context) (*Transaction, error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	return &Transaction{tx}, nil
}

// BeginTx - pgx.BeginTx
func (c *Connection) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (*Transaction, error) {
	tx, err := c.pool.BeginTx(ctx, txOptions)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	return &Transaction{tx}, nil
}

// SendBatch - pgx.SendBatch
func (c *Connection) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return c.pool.SendBatch(ctx, b)
}

// CopyFrom - pgx.CopyFrom
func (c *Connection) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return c.pool.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

// Sqlizer - something that can build sql query
type Sqlizer interface {
	ToSql() (sql string, args []interface{}, err error)
}

// Getx - aka QueryRow
func (c *Connection) Getx(ctx context.Context, dest interface{}, sqlizer Sqlizer) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return fmt.Errorf("postgres: to sql: %w", err)
	}

	return pgxscan.Get(ctx, c.pool, dest, query, args...)
}

// Selectx - aka Query
func (c *Connection) Selectx(ctx context.Context, dest interface{}, sqlizer Sqlizer) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return fmt.Errorf("postgres: to sql: %w", err)
	}

	return pgxscan.Select(ctx, c.pool, dest, query, args...)
}

// Execx - aka Exec
func (c *Connection) Execx(ctx context.Context, sqlizer Sqlizer) (pgconn.CommandTag, error) {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("postgres: to sql: %w", err)
	}

	return c.pool.Exec(ctx, query, args...)
}
