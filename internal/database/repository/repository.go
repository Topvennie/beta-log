// Package repository interacts with the databank and returns models
package repository

import (
	"context"
	"sync"

	"github.com/Topvennie/beta-log/pkg/db"
	"github.com/Topvennie/beta-log/pkg/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	instance *client
	once     sync.Once
)

type client struct {
	db db.DB
}

func Init(db db.DB) {
	once.Do(func() { instance = &client{db: db} })
}

type contextKey string

const queryKey = contextKey("queries")

func queries(ctx context.Context) *sqlc.Queries {
	if q, ok := ctx.Value(queryKey).(*sqlc.Queries); ok {
		return q
	}

	return instance.db.Queries()
}

func WithRollback(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(queryKey).(*sqlc.Queries); ok {
		return fn(ctx)
	}

	return instance.db.WithRollback(ctx, func(q *sqlc.Queries) error {
		txCtx := context.WithValue(ctx, queryKey, q)
		return fn(txCtx)
	})
}

func Pool() *pgxpool.Pool {
	return instance.db.Pool()
}
