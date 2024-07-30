package grpc

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/proto"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

func (s *server) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	var response pb.RegisterUserResponse

	userPassword := ""
	if in.User.Password != nil {
		userPassword = *in.User.Password
	}

	// Validate incoming data
	if in.User.Login == "" {
		return nil, status.Errorf(codes.InvalidArgument, "User login is required")
	}
	if userPassword == "" {
		return nil, status.Errorf(codes.InvalidArgument, "User password is required")
	}

	usr := domain.User{
		Login:    in.User.Login,
		Password: userPassword,
	}

	user, err := s.userService.RegisterUser(ctx, usr)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			slog.Info(err.Error(), "login", in.User.Login)
			return nil, status.Error(codes.AlreadyExists, "This login is already taken")
		} else {
			slog.Error(err.Error())
			return nil, status.Error(codes.Internal, "Something went wrong")
		}
	}

	response.User = &pb.User{
		Id:       user.ID.String(),
		Login:    user.Login,
		Password: nil,
		CreatedAt: &timestamppb.Timestamp{
			Seconds: user.CreatedAt.Unix(),
			Nanos:   int32(user.CreatedAt.Nanosecond()),
		},
	}

	token, err := s.authService.CreateToken(ctx, *user)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "Something went wrong")
	}

	response.Token = token

	return &response, nil
}
