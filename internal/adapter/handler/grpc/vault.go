package grpc

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

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

func getUserIDFromContext(ctx context.Context) (userID uuid.UUID, err error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("X-User-Id")
		if len(values) > 0 {
			err = userID.UnmarshalText([]byte(values[0]))
			return
		}
	}
	err = domain.ErrUserIDNotSet
	return
}

func ProtoVaultRecordToDomain(r *pb.Record) (*domain.VaultRecord, error) {
	if r == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Vault record is required")
	}

	if r.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Vault record id is required")
	}

	recordID, err := uuid.Parse(r.GetId())
	if err != nil {
		slog.Warn(err.Error())
		return nil, status.Errorf(codes.InvalidArgument, "Invalid vault record id")
	}

	if r.GetKind().String() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Vault kind is required")
	}

	recordKind := domain.VaultKind(strings.ToLower(r.GetKind().String()))
	if err := recordKind.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid vault kind")
	}

	recordCreatedAt := time.Now()
	if r.GetCreatedAt() != nil {
		recordCreatedAt = r.GetCreatedAt().AsTime()
	}

	recordUpdatedAt := time.Now()
	if r.GetUpdatedAt() != nil {
		recordUpdatedAt = r.GetUpdatedAt().AsTime()
	}

	var recData domain.IVaultRecordData
	switch r.GetKind() {
	case pb.VaultKind_CREDS:
		if r.GetCreds() == nil {
			return nil, status.Errorf(codes.InvalidArgument, "Creds data is required")
		}

		recData = domain.VaultDataCreds{
			Info:     r.GetInfo(),
			Login:    r.GetCreds().GetLogin(),
			Password: r.GetCreds().GetPassword(),
		}
	case pb.VaultKind_NOTE:
		if r.GetNote() == nil {
			return nil, status.Errorf(codes.InvalidArgument, "Note data is required")
		}

		recData = domain.VaultDataNote{
			Info:    r.GetInfo(),
			Content: r.GetNote().GetContent(),
		}
	case pb.VaultKind_CARD:
		if r.GetCard() == nil {
			return nil, status.Errorf(codes.InvalidArgument, "Card data is required")
		}

		cardExpDate := domain.CardExpDate{}
		if r.GetCard().GetExpDate() != nil {
			cardExpDate = domain.CardExpDate{
				Month: uint(r.GetCard().GetExpDate().GetMonth()),
				Year:  uint(r.GetCard().GetExpDate().GetYear()),
			}
		}

		recData = domain.VaultDataCard{
			Info:       r.GetInfo(),
			Number:     r.GetCard().GetNumber(),
			HolderName: r.GetCard().GetHolderName(),
			ExpDate:    cardExpDate,
			CvcCode:    r.GetCard().GetCvcCode(),
		}
	case pb.VaultKind_FILE:
		if r.GetFile() == nil {
			return nil, status.Errorf(codes.InvalidArgument, "File data is required")
		}

		recData = domain.VaultDataFile{
			Info:     r.GetInfo(),
			Filename: r.GetFile().GetFilename(),
			MimeType: r.GetFile().GetMimeType(),
			Filesize: r.GetFile().GetFilesize(),
			Data:     r.GetFile().GetData(),
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "Undefined vault kind")
	}

	rec := &domain.VaultRecord{
		ID:        recordID,
		Kind:      recordKind,
		Data:      recData,
		CreatedAt: recordCreatedAt,
		UpdatedAt: recordUpdatedAt,
	}

	return rec, nil
}

func DomainVaultRecordToProto(rec domain.VaultRecord) (*pb.Record, error) {
	err := rec.Kind.Validate()
	if err != nil {
		return nil, err
	}

	pbRecord := &pb.Record{
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
		pbRecord.Kind = pb.VaultKind_CREDS
		data := rec.Data.(domain.VaultDataCreds)
		pbRecord.Data = &pb.Record_Creds{
			Creds: &pb.Creds{
				Login:    data.Login,
				Password: data.Password,
			},
		}
		pbRecord.Info = data.Info
	case domain.VaultKindCard:
		pbRecord.Kind = pb.VaultKind_CARD
		data := rec.Data.(domain.VaultDataCard)
		pbRecord.Data = &pb.Record_Card{
			Card: &pb.Card{
				Number:     data.Number,
				HolderName: data.HolderName,
				ExpDate: &pb.CardExpDate{
					Month: uint32(data.ExpDate.Month),
					Year:  uint32(data.ExpDate.Year),
				},
				CvcCode: data.CvcCode,
			},
		}
		pbRecord.Info = data.Info
	case domain.VaultKindNote:
		pbRecord.Kind = pb.VaultKind_NOTE
		data := rec.Data.(domain.VaultDataNote)
		pbRecord.Data = &pb.Record_Note{
			Note: &pb.Note{
				Content: data.Content,
			},
		}
		pbRecord.Info = data.Info
	case domain.VaultKindFile:
		pbRecord.Kind = pb.VaultKind_FILE
		data := rec.Data.(domain.VaultDataFile)
		pbRecord.Data = &pb.Record_File{
			File: &pb.File{
				Filename: data.Filename,
				MimeType: data.MimeType,
				Filesize: data.Filesize,
				Data:     data.Data,
			},
		}
		pbRecord.Info = data.Info
	}

	return pbRecord, nil
}
