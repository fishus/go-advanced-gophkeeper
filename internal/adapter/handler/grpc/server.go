package grpc

import (
	"crypto/tls"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/grpc/interceptor"
	pb "github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/proto"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/port"
)

type Config struct {
	Address     string
	TLSCertFile string
	TLSKeyFile  string
}

type server struct {
	pb.UnimplementedVaultServer
	config       Config
	tokenAdapter port.TokenAdapter
	userService  port.UserService
	authService  port.AuthService
	vaultService port.VaultService
}

func NewServer(
	config Config,
	tokenAdapter port.TokenAdapter,
	userService port.UserService,
	authService port.AuthService,
	vaultService port.VaultService) *server {
	server := server{
		config:       config,
		tokenAdapter: tokenAdapter,
		userService:  userService,
		authService:  authService,
		vaultService: vaultService,
	}

	return &server
}

func (s *server) Serve() error {
	listen, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return fmt.Errorf("error open connection: %w", err)
	}

	// Load TLS certificate and private key
	cert, err := tls.LoadX509KeyPair(s.config.TLSCertFile, s.config.TLSKeyFile)
	if err != nil {
		return fmt.Errorf("error load TLS cert and pricate key: %w", err)
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
	})

	serverOpts := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(
			interceptor.AuthUnaryServerInterceptor(s.tokenAdapter),
		),
	}
	srv := grpc.NewServer(serverOpts...)

	pb.RegisterVaultServer(srv, s)

	if err := srv.Serve(listen); err != nil {
		return fmt.Errorf("error starting the GRPC server: %w", err)
	}

	return nil
}
