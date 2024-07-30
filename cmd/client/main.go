package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	api "github.com/fishus/go-advanced-gophkeeper/internal/adapter/api/grpc"
	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/cli/cobra"
	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/config"
	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/logger"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/service"
)

var buildDate string
var buildVersion string
var secretKey string

func main() {
	rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		err = fmt.Errorf("failed to get path to application directory: %w", err)
		slog.Error(err.Error())
		panic(err)
	}
	cfgFile := filepath.Join(rootDir, "client-config.yaml")
	config, err := config.New(cfgFile)
	if err != nil {
		panic(err)
	}
	config.App.SecretKey = secretKey

	// Set logger
	logger.Set(config.App)

	apiAdapter, err := api.New(config.GRPC.Address)
	if err != nil {
		err = fmt.Errorf("error initializing grpc adapter: %w", err)
		slog.Error(err.Error())
		panic(err)
	}

	clientService := service.NewClientService(apiAdapter)

	cli := cobra.New(clientService)
	_ = cli.Execute(context.Background(), buildDate, buildVersion)
}
