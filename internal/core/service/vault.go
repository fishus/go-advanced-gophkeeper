package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	
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

func (s *vaultService) ListVaultRecords(ctx context.Context, userID uuid.UUID, page, limit uint64) ([]domain.VaultRecord, error) {
	list, err := s.vaultRepo.ListVaultRecords(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *vaultService) GetVaultRecord(ctx context.Context, recID uuid.UUID, userID uuid.UUID) (*domain.VaultRecord, error) {
	rec, err := s.vaultRepo.GetVaultRecordByID(ctx, recID)
	if err != nil {
		return nil, fmt.Errorf("error getting vault record: %w", err)
	}

	if rec.UserID != userID {
		return nil, domain.ErrNotFound
	}

	// NB! Для получение содержимого файла использовать метод GetVaultFileContent
	if rec.Kind == domain.VaultKindFile {
		recData, ok := rec.Data.(domain.VaultDataFile)
		if !ok {
			return nil, domain.ErrInvalidVaultRecordKind
		}
		recData.Data = []byte{}
		rec.Data = recData
	}

	return rec, nil
}

func (s *vaultService) GetVaultFileContent(ctx context.Context, fileID uuid.UUID) ([]byte, error) {
	rec, err := s.vaultRepo.GetVaultRecordByID(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("error getting vault record: %w", err)
	}

	recData, ok := rec.Data.(domain.VaultDataFile)
	if !ok {
		return nil, domain.ErrInvalidVaultRecordKind
	}

	return recData.Data, nil
}
