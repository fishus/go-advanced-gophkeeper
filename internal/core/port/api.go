package port

import (
	"context"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

type ApiAdapter interface {
	Open() error
	Close() error
	LoginUser(ctx context.Context, login, password string) (token string, err error)
	RegisterUser(context.Context, domain.User) (token string, err error)
}
