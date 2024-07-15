package grpc

import (
	pb "github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/proto"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/port"
)

type server struct {
	pb.UnimplementedVaultServer
	userService  port.UserService
	authService  port.AuthService
	vaultService port.VaultService
}

func NewServer(
	userService port.UserService,
	authService port.AuthService,
	vaultService port.VaultService) *server {
	server := server{
		userService:  userService,
		authService:  authService,
		vaultService: vaultService,
	}

	return &server
}
