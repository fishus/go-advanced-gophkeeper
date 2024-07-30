package grpc

import (
	"context"
	"errors"

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
