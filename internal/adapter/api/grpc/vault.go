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

		now := time.Now()

		recordCreatedAt := now
		if v.GetCreatedAt() != nil {
			recordCreatedAt = v.GetCreatedAt().AsTime()
		}

		recordUpdatedAt := now
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

func (api *ApiAdapter) VaultGetRecord(ctx context.Context, recID uuid.UUID) (*domain.VaultRecord, error) {
	resp, err := api.client.GetVaultRecord(ctx, &pb.GetVaultRecordRequest{
		Id: recID.String(),
	})

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.DeadlineExceeded:
				err = domain.ErrTimeout
			case codes.InvalidArgument:
				err = domain.ErrInvalidArgument
			case codes.NotFound:
				err = domain.ErrNotFound
			default:
				err = errors.New(e.Message())
			}
		}
		return nil, err
	}

	record, err := grpc.ProtoVaultRecordToDomain(resp.GetRecord())
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (api *ApiAdapter) VaultGetFile(ctx context.Context, recID uuid.UUID) (*domain.VaultDataFile, []byte, error) {
	resp, err := api.client.DownloadVaultFile(ctx, &pb.DownloadVaultFileRequest{
		Id: recID.String(),
	})

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.DeadlineExceeded:
				err = domain.ErrTimeout
			case codes.InvalidArgument:
				err = domain.ErrInvalidArgument
			case codes.NotFound:
				err = domain.ErrNotFound
			default:
				err = errors.New(e.Message())
			}
		}
		return nil, nil, err
	}

	if resp.GetFile() == nil {
		return nil, nil, domain.ErrInvalidServerResponse
	}

	file := domain.VaultDataFile{
		Filename: resp.GetFile().GetFilename(),
		MimeType: resp.GetFile().GetMimeType(),
		Filesize: resp.GetFile().GetFilesize(),
	}

	data := resp.GetFile().GetData()

	return &file, data, nil
}
