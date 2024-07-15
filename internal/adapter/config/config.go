package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
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

// New creates a new config instance
func New(file string) (*Config, error) {
	err := godotenv.Load(file)
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	app := &App{
		LogLevel:  os.Getenv("APP_LOG_LEVEL"),
		SecretKey: os.Getenv("APP_SECRET_KEY"),
	}

	db := &DB{
		URI: os.Getenv("DB_URI"),
	}

	grpc := &GRPC{
		Address: os.Getenv("GRPC_ADDRESS"),
	}

	tokenDuration, err := time.ParseDuration(os.Getenv("TOKEN_DURATION"))
	if err != nil {
		return nil, fmt.Errorf("parse .env error: invalid token duration value. %w", err)
	}
	token := &Token{
		Duration: tokenDuration,
	}

	return &Config{
		app,
		db,
		grpc,
		token,
	}, nil
}
