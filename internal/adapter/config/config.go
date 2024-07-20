package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type (
	// Config contains environment variables for the application, database, etc
	Config struct {
		App   *App
		DB    *DB
		GRPC  *GRPC
		Token *Token
	}

	// App contains all the environment variables for the application
	App struct {
		LogLevel  string
		SecretKey string
	}

	// DB contains all the environment variables for the database
	DB struct {
		URI string
	}

	// GRPC contains all the environment variables for the grpc server
	GRPC struct {
		Address string
	}

	// Token contains all the environment variables for the token service
	Token struct {
		Duration time.Duration
	}
)

// New creates a new config instance from .env or .yaml file
func New(filename string) (*Config, error) {
	viper.SetConfigFile(filename)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	viper.RegisterAlias("APP_LOG_LEVEL", "app.log_level")
	viper.RegisterAlias("APP_SECRET_KEY", "app.secret_key")
	viper.RegisterAlias("DB_URI", "db.uri")
	viper.RegisterAlias("GRPC_ADDRESS", "grpc.address")
	viper.RegisterAlias("TOKEN_DURATION", "token.duration")

	viper.SetDefault("APP_LOG_LEVEL", "debug")
	viper.SetDefault("GRPC_ADDRESS", ":8080")
	viper.SetDefault("TOKEN_DURATION", "15m")

	app := &App{
		LogLevel:  viper.GetString("APP_LOG_LEVEL"),
		SecretKey: viper.GetString("APP_SECRET_KEY"),
	}

	db := &DB{
		URI: viper.GetString("DB_URI"),
	}

	grpc := &GRPC{
		Address: viper.GetString("GRPC_ADDRESS"),
	}

	tokenDuration := viper.GetDuration("TOKEN_DURATION")
	token := &Token{
		Duration: tokenDuration,
	}

	viper.Reset()

	return &Config{
		app,
		db,
		grpc,
		token,
	}, nil
}
