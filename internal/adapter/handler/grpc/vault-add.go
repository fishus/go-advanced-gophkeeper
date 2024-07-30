package grpc

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/proto"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

func (s *server) AddVaultRecord(ctx context.Context, in *pb.AddVaultRecordRequest) (*pb.AddVaultRecordResponse, error) {
	var response pb.AddVaultRecordResponse

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "Something went wrong")
	}

	rec, err := ProtoVaultRecordToDomain(in.GetRecord())
	if err != nil {
		return nil, err
	}
	rec.UserID = userID

	rec, err = s.vaultService.AddVaultRecord(ctx, *rec)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			slog.Info(err.Error(), "id", rec.ID)
			return nil, status.Error(codes.AlreadyExists, "Vault record already exists")
		} else {
			slog.Error(err.Error())
			return nil, status.Error(codes.Internal, "Something went wrong")
		}
	}

	pbRecord, err := DomainVaultRecordToProto(*rec)
	if err != nil {
		return nil, err
	}

	response.Record = pbRecord

	return &response, nil
}
