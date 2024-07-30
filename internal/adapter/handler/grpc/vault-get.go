package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/proto"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

func (s *server) GetVaultRecord(ctx context.Context, in *pb.GetVaultRecordRequest) (*pb.GetVaultRecordResponse, error) {
	var response pb.GetVaultRecordResponse

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "Something went wrong")
	}

	recordID, err := uuid.Parse(in.GetId())
	if err != nil {
		slog.Warn(err.Error())
		return nil, status.Errorf(codes.InvalidArgument, "Invalid vault record id")
	}

	rec, err := s.vaultService.GetVaultRecord(ctx, recordID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "Vault record not found")
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

func (s *server) DownloadVaultFile(ctx context.Context, in *pb.DownloadVaultFileRequest) (*pb.DownloadVaultFileResponse, error) {
	var response pb.DownloadVaultFileResponse

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "Something went wrong")
	}

	recordID, err := uuid.Parse(in.GetId())
	if err != nil {
		slog.Warn(err.Error())
		return nil, status.Errorf(codes.InvalidArgument, "Invalid vault record id")
	}

	rec, err := s.vaultService.GetVaultRecord(ctx, recordID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "Vault record not found")
		} else {
			slog.Error(err.Error())
			return nil, status.Error(codes.Internal, "Something went wrong")
		}
	}

	recData, ok := rec.Data.(domain.VaultDataFile)
	if !ok {
		slog.Error(domain.ErrInvalidVaultRecordKind.Error())
		return nil, status.Error(codes.Internal, domain.ErrInvalidVaultRecordKind.Error())
	}

	fileContent, err := s.vaultService.GetVaultFileContent(ctx, recordID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "File not found")
		} else {
			slog.Error(err.Error())
			return nil, status.Error(codes.Internal, "Something went wrong")
		}
	}

	response.File = &pb.File{
		Filename: recData.Filename,
		MimeType: recData.MimeType,
		Filesize: recData.Filesize,
		Data:     fileContent,
	}

	return &response, nil
}
