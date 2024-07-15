package domain

import (
	"time"

	"github.com/google/uuid"
)

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

/*
	x := make([]domain.IRecord, 0, 2)
	x = append(x, recA)
	x = append(x, recB)

	for _, val := range x {

		if xa, ok := val.(domain.IRecordA); ok {
			fmt.Printf("xa: %#v\n", xa.GetA())
		}

		if xb, ok := val.(domain.IRecordB); ok {
			fmt.Printf("xb: %#v\n", xb.GetB())
		}
	}
*/
