package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

type ApiAdapter interface {
	Open() error
	Close() error
	LoginUser(ctx context.Context, login, password string) (token string, err error)
	RegisterUser(context.Context, domain.User) (token string, err error)
	SetToken(ctx context.Context, token string) (context.Context, error)
	VaultAddRecord(context.Context, domain.VaultRecord) (*domain.VaultRecord, error)
	VaultListRecords(ctx context.Context, page, limit uint64) ([]domain.VaultListItem, error)
	VaultGetRecord(ctx context.Context, recID uuid.UUID) (*domain.VaultRecord, error)
	VaultGetFile(ctx context.Context, recID uuid.UUID) (*domain.VaultDataFile, []byte, error)
}
