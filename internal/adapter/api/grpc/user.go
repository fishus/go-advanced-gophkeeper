package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/proto"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

func (api *ApiAdapter) LoginUser(ctx context.Context, login, password string) (token string, err error) {
	resp, err := api.client.LoginUser(ctx, &pb.LoginUserRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				err = domain.ErrUserNotFound
			case codes.DeadlineExceeded:
				err = domain.ErrTimeout
			case codes.InvalidArgument, codes.Unauthenticated:
				err = domain.ErrInvalidCredentials
			default:
				err = errors.New(e.Message())
			}
		}
		return
	}
	return resp.Token, nil
}

func (api *ApiAdapter) RegisterUser(ctx context.Context, login, password string) (token string, err error) {
	createdAt := time.Now()
	resp, err := api.client.RegisterUser(ctx, &pb.RegisterUserRequest{
		User: &pb.User{
			Id:       uuid.NewString(),
			Login:    login,
			Password: &password,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: createdAt.Unix(),
				Nanos:   int32(createdAt.Nanosecond()),
			},
		},
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
		return
	}
	return resp.Token, nil
}
