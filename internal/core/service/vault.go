package service

import (
	"context"
	"fmt"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/port"
)

type vaultService struct {
	vaultRepo port.VaultRepository
}

func NewVaultService(vaultRepo port.VaultRepository) port.VaultService {
	return &vaultService{
		vaultRepo: vaultRepo,
	}
}

func (s *vaultService) AddVaultRecord(ctx context.Context, r domain.VaultRecord) (*domain.VaultRecord, error) {
	if err := r.Data.Validate(); err != nil {
		return nil, err
	}

	rec, err := s.vaultRepo.CreateVaultRecord(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("vault record not created: %w", err)
	}

	return rec, nil
}
