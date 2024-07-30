package paseto

import (
	"time"

	"aidanwoods.dev/go-paseto"

	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/config"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/port"
)

// TokenAdapter implements port.TokenAdapter interface
// and provides access to the paseto library
type TokenAdapter struct {
	token    *paseto.Token
	key      *paseto.V4SymmetricKey
	parser   *paseto.Parser
	duration time.Duration
}

// New creates a new paseto instance
func New(config *config.Token) (port.TokenAdapter, error) {
	token := paseto.NewToken()
	key := paseto.NewV4SymmetricKey()
	parser := paseto.NewParser()

	return &TokenAdapter{
		&token,
		&key,
		&parser,
		config.Duration,
	}, nil
}

// CreateToken creates a new paseto token
func (t *TokenAdapter) CreateToken(payload domain.TokenPayload) (string, error) {
	err := t.token.Set("payload", payload)
	if err != nil {
		return "", err
	}

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(t.duration)

	t.token.SetIssuedAt(issuedAt)
	t.token.SetNotBefore(issuedAt)
	t.token.SetExpiration(expiredAt)

	token := t.token.V4Encrypt(*t.key, nil)

	return token, nil
}

// VerifyToken verifies the paseto token
func (t *TokenAdapter) VerifyToken(token string) (*domain.TokenPayload, error) {
	var payload *domain.TokenPayload

	parsedToken, err := t.parser.ParseV4Local(*t.key, token, nil)
	if err != nil {
		if err.Error() == "this token has expired" {
			return nil, domain.ErrExpiredToken
		}
		return nil, domain.ErrInvalidToken
	}

	err = parsedToken.Get("payload", &payload)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	return payload, nil
}
