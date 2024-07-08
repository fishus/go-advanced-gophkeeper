package repository

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/repository/postgres"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

// UserRepository implements port.UserRepository interface
// and provides an access to the postgres database
type UserRepository struct {
	db *postgres.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

// CreateUser creates a new user in the database
func (repo *UserRepository) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	query := repo.db.QueryBuilder.Insert("users").
		Columns("id", "login", "password", "created_at").
		Values(user.ID, user.Login, user.Password, user.CreatedAt).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = repo.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if errCode := repo.db.ErrorCode(err); errCode == pgerrcode.UniqueViolation {
			return nil, domain.ErrAlreadyExists
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByLogin gets a user by login from the database
func (repo *UserRepository) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User

	query := repo.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"login": login}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = repo.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
