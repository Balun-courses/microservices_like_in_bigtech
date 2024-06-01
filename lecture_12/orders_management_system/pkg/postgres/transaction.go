package postgres

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type Transaction struct {
	pgx.Tx
}

func (t *Transaction) Getx(ctx context.Context, dest interface{}, sqlizer Sqlizer) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return err
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "postgres.tx.Getx")
	defer span.Finish()

	span.LogFields(
		log.String("query", query),
		log.Object("args", args),
	)

	return pgxscan.Get(ctx, t.Tx, dest, query, args...)
}

func (t *Transaction) Selectx(ctx context.Context, dest interface{}, sqlizer Sqlizer) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return err
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "postgres.tx.Selectx")
	defer span.Finish()

	span.LogFields(
		log.String("query", query),
		log.Object("args", args),
	)

	return pgxscan.Select(ctx, t.Tx, dest, query, args)
}

func (t *Transaction) Execx(ctx context.Context, sqlizer Sqlizer) (pgconn.CommandTag, error) {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "postgres.tx.Execx")
	defer span.Finish()

	span.LogFields(
		log.String("query", query),
		log.Object("args", args),
	)

	return t.Tx.Exec(ctx, query, args...)
}
