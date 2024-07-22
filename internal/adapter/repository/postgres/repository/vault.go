package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgerrcode"

	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/repository/postgres"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/port"
)

// VaultRepository implements port.VaultRepository interface
// and provides access to the postgres database
type VaultRepository struct {
	db           *postgres.DB
	cryptAdapter port.CryptAdapter
}

// NewVaultRepository creates a new vault repository instance
func NewVaultRepository(db *postgres.DB, cryptAdapter port.CryptAdapter) *VaultRepository {
	return &VaultRepository{
		db,
		cryptAdapter,
	}
}

// CreateVaultRecord creates a new vault record in the database
func (repo *VaultRepository) CreateVaultRecord(ctx context.Context, rec domain.VaultRecord) (*domain.VaultRecord, error) {
	// FIXME

	src, err := json.Marshal(rec.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal vault data error: %w", err)
	}

	data, err := repo.cryptAdapter.EncryptSymmetric(ctx, src)
	if err != nil {
		return nil, fmt.Errorf("encrypt vault data error: %w", err)
	}

	query := repo.db.QueryBuilder.Insert("vault").
		Columns("id", "user_id", "kind", "data", "created_at", "updated_at").
		Values(rec.ID, rec.UserID, rec.Kind, data, rec.CreatedAt, rec.UpdatedAt).
		Suffix("RETURNING id, user_id, kind, created_at, updated_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = repo.db.QueryRow(ctx, sql, args...).Scan(
		&rec.ID,
		&rec.UserID,
		&rec.Kind,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	)
	if err != nil {
		if errCode := repo.db.ErrorCode(err); errCode == pgerrcode.UniqueViolation {
			return nil, domain.ErrAlreadyExists
		}
		return nil, err
	}

	return &rec, nil

	//aaa, err := repo.cryptAdapter.DecryptSymmetric(ctx, dst)

	//var uRecData domain.VaultDataCard
	//if err = json.Unmarshal(out, &uRecData); err != nil {
	//	return nil, err
	//}
	//fmt.Printf("Unmarshal: %#v\n", uRecData)
}

//func (repo *VaultRepository) marshalRecordData(ctx context.Context, rec domain.IVaultRecord) (string, error) {
//	switch rec.GetKind() {
//	case domain.VaultKindCreds:
//
//	default:
//		return "", fmt.Errorf("unknown vault kind")
//	}
//	return "test", nil
//}

//func (repo *VaultRepository) unmarshalRecordData(ctx context.Context, data string) (rec domain.IVaultRecord, error) {
//	fmt.Printf("secret key: %v\n", repo.secretKey)
//	//s := sha256.Sum256(src)
//	return "test", nil
//}

//func (repo *VaultRepository) encodeVaultRecord(ctx context.Context, rec domain.IVaultRecord) (string, error) {
//	fmt.Printf("secret key: %v\n", repo.secretKey)
//	//s := sha256.Sum256(src)
//	return "test", nil
//}
