package port

import (
	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

// TokenAdapter is an interface for interacting with token-related business logic
type TokenAdapter interface {
	// CreateToken creates a new token for a given user
	CreateToken(domain.TokenPayload) (token string, err error)
	// VerifyToken verifies the token and returns the payload
	VerifyToken(token string) (*domain.TokenPayload, error)
}
