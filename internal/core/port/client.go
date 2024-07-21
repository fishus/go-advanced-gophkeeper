package port

import (
	"context"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

type ClientService interface {
	Setup(context.Context) error
	Teardown(context.Context) error
	SetToken(token string) ClientService
	UserLogin(ctx context.Context, login, password string) (token string, err error)
	UserRegister(ctx context.Context, login, password string) (token string, err error)
	VaultAddNote(ctx context.Context, note domain.VaultDataNote) (*domain.VaultRecord, error)
}
