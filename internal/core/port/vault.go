package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

type VaultService interface {
	// AddVaultRecord adds new record into vault
	AddVaultRecord(context.Context, domain.VaultRecord) (*domain.VaultRecord, error)
	ListVaultRecords(ctx context.Context, userID uuid.UUID, page, limit uint64) ([]domain.VaultRecord, error)
}

type VaultRepository interface {
	// CreateVaultRecord inserts a new vault's record into the database
	CreateVaultRecord(context.Context, domain.VaultRecord) (*domain.VaultRecord, error)
	ListVaultRecords(ctx context.Context, userID uuid.UUID, page, limit uint64) ([]domain.VaultRecord, error)
}
