package domain

import (
	"time"

	"github.com/google/uuid"
)

type VaultListItem struct {
	ID        uuid.UUID
	Kind      VaultKind
	Info      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type VaultRecord struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Kind      VaultKind
	Data      IVaultRecordData
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IVaultRecordData interface {
	GetInfo() string
	Validate() error
}

// VaultKind Enumeration of record types in storage
type VaultKind string

const (
	VaultKindUndefined VaultKind = ""
	VaultKindCreds     VaultKind = "creds"
	VaultKindNote      VaultKind = "note"
	VaultKindCard      VaultKind = "card"
	VaultKindFile      VaultKind = "file"
)

func (k VaultKind) Validate() error {
	switch k {
	case VaultKindCreds,
		VaultKindNote,
		VaultKindCard,
		VaultKindFile:
		return nil
	case VaultKindUndefined:
		return ErrUndefinedVaultKind
	}
	return ErrInvalidVaultKind
}

func (k VaultKind) String() string {
	return string(k)
}
