package grpc

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/grpc"
	pb "github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/proto"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

func (api *ApiAdapter) VaultAddRecord(ctx context.Context, r domain.VaultRecord) (*domain.VaultRecord, error) {
	pbRec, err := grpc.DomainVaultRecordToProto(r)
	if err != nil {
		return nil, err
	}

	resp, err := api.client.AddVaultRecord(ctx, &pb.AddVaultRecordRequest{
		Record: pbRec,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.DeadlineExceeded:
				err = domain.ErrTimeout
			case codes.AlreadyExists:
				err = domain.ErrAlreadyExists
			case codes.InvalidArgument:
				err = domain.ErrInvalidArgument
			default:
				err = errors.New(e.Message())
			}
		}
		return nil, err
	}

	rec, err := grpc.ProtoVaultRecordToDomain(resp.GetRecord())
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func (api *ApiAdapter) VaultListRecords(ctx context.Context, page, limit uint64) ([]domain.VaultListItem, error) {
	resp, err := api.client.ListVaultRecords(ctx, &pb.ListVaultRecordsRequest{
		Page:  page,
		Limit: limit,
	})

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.DeadlineExceeded:
				err = domain.ErrTimeout
			case codes.InvalidArgument:
				err = domain.ErrInvalidArgument
			default:
				err = errors.New(e.Message())
			}
		}
		return nil, err
	}

	var list []domain.VaultListItem
	for _, v := range resp.GetList() {
		recordID, err := uuid.Parse(v.GetId())
		if err != nil {
			return nil, err
		}

		recordKind := domain.VaultKind(strings.ToLower(v.GetKind().String()))
		if err := recordKind.Validate(); err != nil {
			return nil, err
		}

		recordCreatedAt := time.Now()
		if v.GetCreatedAt() != nil {
			recordCreatedAt = v.GetCreatedAt().AsTime()
		}

		recordUpdatedAt := time.Now()
		if v.GetUpdatedAt() != nil {
			recordUpdatedAt = v.GetUpdatedAt().AsTime()
		}

		item := domain.VaultListItem{
			ID:        recordID,
			Kind:      recordKind,
			Info:      v.Info,
			CreatedAt: recordCreatedAt,
			UpdatedAt: recordUpdatedAt,
		}

		list = append(list, item)
	}

	return list, nil
}
