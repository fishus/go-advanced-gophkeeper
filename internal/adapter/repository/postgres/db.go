package postgres

import (
	"context"
	"embed"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	
	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/config"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// DB is a wrapper for PostgreSQL database connection
// that uses pgxpool as database driver.
// It also holds a reference to squirrel.StatementBuilderType
// which is used to build SQL queries that compatible with PostgreSQL syntax
type DB struct {
	*pgxpool.Pool
	Config       config.DB
	QueryBuilder *squirrel.StatementBuilderType
}

// New creates a new PostgreSQL database instance
func New(ctx context.Context, config config.DB) (*DB, error) {
	db, err := pgxpool.New(ctx, config.URI)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &DB{
		db,
		config,
		&psql,
	}, nil
}

// Migrate runs the database migration
func (db *DB) Migrate() error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	conn := stdlib.OpenDBFromPool(db.Pool)
	defer conn.Close()

	err := goose.Up(conn, "migrations")

	if err != nil {
		return err
	}

	return nil
}

// ErrorCode returns the error code of the given error
func (db *DB) ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	errors.As(err, &pgErr)
	return pgErr.Code
}

// Close closes the database connection
func (db *DB) Close() error {
	db.Pool.Close()
	return nil
}
