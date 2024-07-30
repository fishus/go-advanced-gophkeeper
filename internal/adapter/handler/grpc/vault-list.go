package grpc

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/proto"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

func (s *server) ListVaultRecords(ctx context.Context, in *pb.ListVaultRecordsRequest) (*pb.ListVaultRecordsResponse, error) {
	var response pb.ListVaultRecordsResponse

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "Something went wrong")
	}

	page := in.GetPage()
	limit := in.GetLimit()

	list, err := s.vaultService.ListVaultRecords(ctx, userID, page, limit)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "Something went wrong")
	}

	for _, rec := range list {
		pbItem := pb.ListVaultRecordsResponse_ListItem{
			Id: rec.ID.String(),
			CreatedAt: &timestamppb.Timestamp{
				Seconds: rec.CreatedAt.Unix(),
				Nanos:   int32(rec.CreatedAt.Nanosecond()),
			},
			UpdatedAt: &timestamppb.Timestamp{
				Seconds: rec.UpdatedAt.Unix(),
				Nanos:   int32(rec.UpdatedAt.Nanosecond()),
			},
		}
		switch rec.Kind {
		case domain.VaultKindCreds:
			pbItem.Kind = pb.VaultKind_CREDS
			data := rec.Data.(domain.VaultDataCreds)
			pbItem.Info = data.Info
		case domain.VaultKindCard:
			pbItem.Kind = pb.VaultKind_CARD
			data := rec.Data.(domain.VaultDataCard)
			pbItem.Info = data.Info
		case domain.VaultKindNote:
			pbItem.Kind = pb.VaultKind_NOTE
			data := rec.Data.(domain.VaultDataNote)
			pbItem.Info = data.Info
		case domain.VaultKindFile:
			pbItem.Kind = pb.VaultKind_FILE
			data := rec.Data.(domain.VaultDataFile)
			pbItem.Info = data.Info
		default:
			slog.Error(err.Error())
			return nil, status.Error(codes.Internal, "Something went wrong")
		}

		response.List = append(response.List, &pbItem)
	}

	return &response, nil
}
