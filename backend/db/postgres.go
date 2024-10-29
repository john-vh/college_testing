package db

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxConn interface {
	Begin(ctx context.Context) (*PgxQueries, error)
}

type PgxStore struct {
	pool *pgxpool.Pool
}

type PgxQueries struct {
	tx pgx.Tx
}

func NewPgxStorage(ctx context.Context, config *pgxpool.Config) (*PgxStore, error) {
	pool, err := pgxpool.New(ctx, config.ConnString())
	if err != nil {
		return nil, err
	}

	pg := &PgxStore{
		pool: pool,
	}
	return pg, nil
}

func (ps *PgxStore) Begin(ctx context.Context) (*PgxQueries, error) {
	tx, err := ps.pool.Begin(ctx)
	if err != nil {
		return nil, handlePgxError(err)
	}

	return &PgxQueries{tx: tx}, nil
}

func (pq *PgxQueries) Begin(ctx context.Context) (*PgxQueries, error) {
	tx, err := pq.tx.Begin(ctx)
	if err != nil {
		return nil, handlePgxError(err)
	}
	return &PgxQueries{tx: tx}, nil
}

func (pq *PgxQueries) Rollback(ctx context.Context) error {
	return pq.tx.Rollback(ctx)
}

func (pq *PgxQueries) Commit(ctx context.Context) error {
	return pq.tx.Commit(ctx)
}

func WithTx(ctx context.Context, conn PgxConn, fn func(*PgxQueries) error) error {
	pq, err := conn.Begin(ctx)
	if err != nil {
		return handlePgxError(err)
	}
	defer pq.Rollback(ctx)

	err = fn(pq)
	if err != nil {
		return err
	}

	err = pq.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func WithTxRet[T any](ctx context.Context, conn PgxConn, fn func(*PgxQueries) (T, error)) (T, error) {
	pq, err := conn.Begin(ctx)
	if err != nil {
		return *new(T), handlePgxError(err)
	}
	defer pq.Rollback(ctx)

	val, err := fn(pq)
	if err != nil {
		return val, err
	}

	err = pq.Commit(ctx)
	if err != nil {
		return val, handlePgxError(err)
	}

	return val, nil
}

func handlePgxError(err error) error {
	var pgerr *pgconn.PgError
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNoRows
	}

	if errors.As(err, &pgerr) {
		switch pgerr.Code {
		case pgerrcode.ForeignKeyViolation, pgerrcode.InvalidForeignKey:
			return ErrForeignKey
		case pgerrcode.UniqueViolation:
			return ErrUnique
		}
	}

	return ErrDB
}
