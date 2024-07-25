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
	VaultAddCard(ctx context.Context, card domain.VaultDataCard) (*domain.VaultRecord, error)
	VaultAddCreds(ctx context.Context, creds domain.VaultDataCreds) (*domain.VaultRecord, error)
	VaultAddFile(ctx context.Context, file domain.VaultDataFile) (*domain.VaultRecord, error)
	VaultAddRecord(ctx context.Context, data domain.IVaultRecordData) (*domain.VaultRecord, error)
}
