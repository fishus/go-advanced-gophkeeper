package repository

import (
	"context"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
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

func (repo *VaultRepository) ListVaultRecords(ctx context.Context, userID uuid.UUID, page, limit uint64) ([]domain.VaultRecord, error) {
	query := repo.db.QueryBuilder.Select("id, kind, data, created_at, updated_at").
		From("vault").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		Limit(limit).
		Offset((page - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var encodedData []byte
	var record domain.VaultRecord
	var records []domain.VaultRecord

	for rows.Next() {
		err = rows.Scan(
			&record.ID,
			&record.Kind,
			&encodedData,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		jsonData, err := repo.cryptAdapter.DecryptSymmetric(ctx, encodedData)

		switch record.Kind {
		case domain.VaultKindNote:
			var recData domain.VaultDataNote
			if err = json.Unmarshal(jsonData, &recData); err != nil {
				return nil, err
			}
			record.Data = recData
		case domain.VaultKindCard:
			var recData domain.VaultDataCard
			if err = json.Unmarshal(jsonData, &recData); err != nil {
				return nil, err
			}
			record.Data = recData
		case domain.VaultKindCreds:
			var recData domain.VaultDataCreds
			if err = json.Unmarshal(jsonData, &recData); err != nil {
				return nil, err
			}
			record.Data = recData
		case domain.VaultKindFile:
			var recData domain.VaultDataFile
			if err = json.Unmarshal(jsonData, &recData); err != nil {
				return nil, err
			}
			recData.Data = []byte{}
			record.Data = recData
		default:
			return nil, domain.ErrInvalidVaultKind
		}

		records = append(records, record)
	}

	return records, nil
}
