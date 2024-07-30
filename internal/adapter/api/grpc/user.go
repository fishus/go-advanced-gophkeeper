package grpc

import (
	"context"
	"errors"

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

func (api *ApiAdapter) RegisterUser(ctx context.Context, user domain.User) (token string, err error) {
	resp, err := api.client.RegisterUser(ctx, &pb.RegisterUserRequest{
		User: &pb.User{
			Id:       user.ID.String(),
			Login:    user.Login,
			Password: &user.Password,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: user.CreatedAt.Unix(),
				Nanos:   int32(user.CreatedAt.Nanosecond()),
			},
		},
	})
	if err != nil {
		err = handleErrCodes(err)
		return
	}
	return resp.Token, nil
}
