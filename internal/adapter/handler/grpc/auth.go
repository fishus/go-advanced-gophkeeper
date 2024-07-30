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

func (s *server) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	var response pb.LoginUserResponse

	// Validate incoming data
	if in.Login == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Login cannot be empty")
	}
	if in.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Password cannot be empty")
	}

	token, err := s.authService.LoginUser(ctx, in.Login, in.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "Incorrect login or password")
		} else {
			slog.Error(err.Error())
			return nil, status.Error(codes.Internal, "Something went wrong")
		}
	}

	response.Token = token

	return &response, nil
}
