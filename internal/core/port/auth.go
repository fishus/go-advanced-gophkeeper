package port

import (
	"context"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

// AuthService is an interface for interacting with user authentication-related business logic
type AuthService interface {
	// LoginUser authenticates a user by login and password and returns a token
	LoginUser(ctx context.Context, login, password string) (token string, err error)
	CreateToken(context.Context, domain.User) (token string, err error)
}
