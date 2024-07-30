package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/auth/paseto"
	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/config"
	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/crypt/gsm"
	handler "github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/grpc"
	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/logger"
	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/repository/postgres"
	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/repository/postgres/repository"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/service"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/util"
)

var buildDate string
var buildVersion string

func main() {
	fmt.Println(util.GetBuildInfo(buildDate, buildVersion))

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

	// Crypto adapter
	cryptAdapter, err := gsm.New([]byte(config.App.SecretKey))
	if err != nil {
		err = fmt.Errorf("error initializing crypto adapter: %w", err)
		slog.Error(err.Error())
		panic(err)
	}

	// Vault repository adapter and vault service
	vaultRepo := repository.NewVaultRepository(db, cryptAdapter)
	vaultService := service.NewVaultService(vaultRepo)

	fmt.Println("Starting the grpc server")

	// Handler adapter
	server := handler.NewServer(
		handler.Config{
			Address:     config.GRPC.Address,
			TLSCertFile: config.TLS.CertFile,
			TLSKeyFile:  config.TLS.KeyFile,
		}, tokenAdapter, userService, authService, vaultService)
	err = server.Serve()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}
