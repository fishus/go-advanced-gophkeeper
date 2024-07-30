package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/port"
)

type clientService struct {
	apiAdapter port.ApiAdapter
	token      string
}

func NewClientService(
	apiAdapter port.ApiAdapter,
) port.ClientService {
	return &clientService{
		apiAdapter: apiAdapter,
	}
}

func (s *clientService) Setup(ctx context.Context) error {
	return s.apiAdapter.Open()
}

func (s *clientService) Teardown(ctx context.Context) error {
	return s.apiAdapter.Close()
}

func (s *clientService) SetToken(token string) port.ClientService {
	s.token = token
	return s
}

func (s *clientService) UserLogin(ctx context.Context, login, password string) (token string, err error) {
	token, err = s.apiAdapter.LoginUser(ctx, login, password)
	if err == nil {
		s.token = token
	}
	return
}

func (s *clientService) UserRegister(ctx context.Context, login, password string) (token string, err error) {
	user := domain.User{
		ID:        uuid.New(),
		Login:     login,
		Password:  password,
		CreatedAt: time.Now(),
	}
	token, err = s.apiAdapter.RegisterUser(ctx, user)
	if err == nil {
		s.token = token
	}
	return
}

func (s *clientService) VaultAddNote(ctx context.Context, data domain.VaultDataNote) (*domain.VaultRecord, error) {
	return s.VaultAddRecord(ctx, data)
}

func (s *clientService) VaultAddCard(ctx context.Context, data domain.VaultDataCard) (*domain.VaultRecord, error) {
	return s.VaultAddRecord(ctx, data)
}

func (s *clientService) VaultAddCreds(ctx context.Context, data domain.VaultDataCreds) (*domain.VaultRecord, error) {
	return s.VaultAddRecord(ctx, data)
}

func (s *clientService) VaultAddFile(ctx context.Context, data domain.VaultDataFile) (*domain.VaultRecord, error) {
	return s.VaultAddRecord(ctx, data)
}

func (s *clientService) VaultAddRecord(ctx context.Context, data domain.IVaultRecordData) (*domain.VaultRecord, error) {
	var kind domain.VaultKind
	switch data.(type) {
	case domain.VaultDataNote:
		kind = domain.VaultKindNote
	case domain.VaultDataCard:
		kind = domain.VaultKindCard
	case domain.VaultDataCreds:
		kind = domain.VaultKindCreds
	case domain.VaultDataFile:
		kind = domain.VaultKindFile
	default:
		return nil, domain.ErrInvalidVaultKind
	}

	if err := data.Validate(); err != nil {
		return nil, err
	}

	now := time.Now()

	rec := &domain.VaultRecord{
		ID:        uuid.New(),
		Kind:      kind,
		Data:      data,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// TODO save into local db

	ctx, err := s.apiAdapter.SetToken(ctx, s.token)
	if err != nil {
		return nil, err
	}
	rec, err = s.apiAdapter.VaultAddRecord(ctx, *rec)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func (s *clientService) VaultListRecords(ctx context.Context, page, limit uint64) ([]domain.VaultListItem, error) {
	ctx, err := s.apiAdapter.SetToken(ctx, s.token)
	if err != nil {
		return nil, err
	}

	// TODO load from local db

	return s.apiAdapter.VaultListRecords(ctx, page, limit)
}

func (s *clientService) VaultGetRecord(ctx context.Context, recID uuid.UUID) (*domain.VaultRecord, error) {
	ctx, err := s.apiAdapter.SetToken(ctx, s.token)
	if err != nil {
		return nil, err
	}

	record, err := s.apiAdapter.VaultGetRecord(ctx, recID)
	if err != nil {
		return nil, err
	}

	// TODO load from local db

	return record, nil
}

func (s *clientService) VaultGetFile(ctx context.Context, recID uuid.UUID) (*domain.VaultDataFile, []byte, error) {
	ctx, err := s.apiAdapter.SetToken(ctx, s.token)
	if err != nil {
		return nil, nil, err
	}

	file, data, err := s.apiAdapter.VaultGetFile(ctx, recID)
	if err != nil {
		return nil, nil, err
	}

	// TODO save file to local storage
	// TODO load from local db

	return file, data, nil
}
