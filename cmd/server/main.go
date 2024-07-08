package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/auth/paseto"
	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/config"
	handler "github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/grpc"
	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/grpc/interceptor"
	pb "github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/proto"
	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/logger"
	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/repository/postgres"
	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/repository/postgres/repository"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/service"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/util"
)

var buildDate string
var buildVersion string

func main() {
	util.PrintBuildInfo(buildDate, buildVersion)

	// Load environment variables
	config, err := config.New(".server.env")
	if err != nil {
		err = fmt.Errorf("error loading environment variables: %w", err)
		slog.Error(err.Error())
		panic(err)
	}

	// Set logger
	logger.Set(config.App)

	slog.Info("Starting the server")

	// Init database adapter
	ctx := context.Background()
	db, err := postgres.New(ctx, *config.DB)
	if err != nil {
		err = fmt.Errorf("error initializing database connection: %w", err)
		slog.Error(err.Error())
		panic(err)
	}
	defer db.Close()

	slog.Info("Successfully connected to the database")

	// Migrate database
	err = db.Migrate()
	if err != nil {
		err = fmt.Errorf("error migrating database: %w", err)
		slog.Error(err.Error())
		panic(err)
	}

	slog.Info("Successfully migrated the database")

	// Dependency injection
	// User repository adapter and user service
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// Init token adapter
	tokenAdapter, err := paseto.New(config.Token)
	if err != nil {
		err = fmt.Errorf("error initializing token adapter: %w", err)
		slog.Error(err.Error())
		panic(err)
	}

	// Auth service
	authService := service.NewAuthService(userRepo, tokenAdapter)

	// Handler adapter
	server := handler.NewServer(userService, authService)

	listen, err := net.Listen("tcp", config.GRPC.Address)
	if err != nil {
		err = fmt.Errorf("error open connection: %w", err)
		slog.Error(err.Error())
		panic(err)
	}

	serverOpts := []grpc.ServerOption{}
	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(
		interceptor.AuthUnaryServerInterceptor(tokenAdapter),
	))
	s := grpc.NewServer(serverOpts...)

	pb.RegisterVaultServer(s, server)

	if err := s.Serve(listen); err != nil {
		err = fmt.Errorf("error starting the GRPC server: %w", err)
		slog.Error(err.Error())
		panic(err)
	}

	slog.Info("Successfully started the GRPC server", "address", config.GRPC.Address)
}
