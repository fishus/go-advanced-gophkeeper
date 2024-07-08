package port

import (
	"context"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

type UserService interface {
	// RegisterUser registers a new user
	RegisterUser(ctx context.Context, user domain.User) (*domain.User, error)
}

type UserRepository interface {
	// CreateUser inserts a new user into the database
	CreateUser(context.Context, domain.User) (*domain.User, error)
	// GetUserByLogin selects a user by login
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
}
