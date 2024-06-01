package transaction_manager

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/moguchev/microservices_courcse/orders_management_system/pkg/postgres"
)

// TransactionManager - менеджер транзакций: позовляет выполнять функции разных репозиториев ходящих в одну БД в рамках транзакции
type TransactionManager struct {
	connection *postgres.Connection
}

// New constructs TransactionManager
func New(connection *postgres.Connection) *TransactionManager {
	return &TransactionManager{connection: connection}
}

type key string

const (
	txKey key = "tx"
)

func (m *TransactionManager) runTransaction(ctx context.Context, txOpts pgx.TxOptions, fn func(ctx context.Context) error) (err error) {
	// If it's nested Transaction, skip initiating a new one and return func(ctx context.Context) error
	tx, ok := ctx.Value(txKey).(*postgres.Transaction)
	if ok {
		return fn(ctx)
	}

	// Begin runTransaction
	pgxTx, err := m.connection.BeginTx(ctx, txOpts)
	if err != nil {
		return fmt.Errorf("can't begin transaction: %v", err)
	}

	tx = &postgres.Transaction{Tx: pgxTx}
	// Set txKey to context
	ctx = context.WithValue(ctx, txKey, tx)

	// Set up a defer function for rolling back the runTransaction.
	defer func() {
		// recover from panic
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}

		// if func(ctx context.Context) error didn't return error - commit
		if err == nil {
			// if commit returns error -> rollback
			err = tx.Commit(ctx)
			if err != nil {
				err = fmt.Errorf("commit failed: %v", err)
			}
		}

		// rollback on any error
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = fmt.Errorf("rollback failed: %v", errRollback)
			}
		}
	}()

	// Execute the code inside the runTransaction. If the function
	// fails, return the error and the defer function will roll back or commit otherwise.

	// return error without wrapping errors.Wrap
	err = fn(ctx)

	return err
}

// GetQueryEngine provides QueryEngine
func (m *TransactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	// Transaction always runs on node with NodeRoleWrite role
	if tx, ok := ctx.Value(txKey).(QueryEngine); ok {
		return tx
	}

	return m.connection
}

// RunReadCommitted execs f func in runTransaction with LevelReadCommitted isolation level
func (m *TransactionManager) RunReadCommitted(ctx context.Context, accessMode pgx.TxAccessMode, f func(ctx context.Context) error) error {
	return m.runTransaction(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: accessMode,
	}, f)
}

// RunRepeatableRead execs f func in runTransaction with LevelRepeatableRead isolation level
func (m *TransactionManager) RunRepeatableRead(ctx context.Context, accessMode pgx.TxAccessMode, f func(ctx context.Context) error) error {
	return m.runTransaction(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: accessMode,
	}, f)
}

// RunSerializable execs f func in runTransaction with LevelSerializable isolation level
func (m *TransactionManager) RunSerializable(ctx context.Context, accessMode pgx.TxAccessMode, f func(ctx context.Context) error) error {
	return m.runTransaction(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: accessMode,
	}, f)
}
