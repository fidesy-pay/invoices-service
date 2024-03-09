package storage

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound      = errors.New("entity not found")
	ErrAlreadyExists = errors.New("entity already exists")
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: pool,
	}
}

func Builder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func Getx[T any](ctx context.Context, executor *pgxpool.Pool, sqlizer sq.Sqlizer) (*T, error) {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return nil, fmt.Errorf("sqlizer.ToSql: %v", err)
	}

	var model T
	err = pgxscan.Get(ctx, executor, &model, query, args...)
	if err != nil {
		return nil, handlePgError(err)
	}

	return &model, nil
}

func Selectx[T any](ctx context.Context, executor *pgxpool.Pool, sqlizer sq.Sqlizer) ([]*T, error) {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return nil, fmt.Errorf("sqlizer.ToSql: %v", err)
	}

	var model []*T
	err = pgxscan.Select(ctx, executor, &model, query, args...)
	if err != nil {
		return nil, handlePgError(err)
	}

	return model, nil
}

func (s *Storage) Exec(ctx context.Context, sqlizer sq.Sqlizer) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return fmt.Errorf("sqlizer.ToSql: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args)
	return err
}

func handlePgError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// Check for PostgreSQL unique constraint violation error
		if pgErr.Code == "23505" {
			return ErrAlreadyExists
		}
	}

	return err
}
