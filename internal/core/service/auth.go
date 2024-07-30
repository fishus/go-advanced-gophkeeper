package service

import (
	"context"
	"errors"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/port"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/util"
)

type authService struct {
	userRepo     port.UserRepository
	tokenAdapter port.TokenAdapter
}

func NewAuthService(userRepo port.UserRepository, tokenAdapter port.TokenAdapter) port.AuthService {
	return &authService{
		userRepo:     userRepo,
		tokenAdapter: tokenAdapter,
	}
}

func (s *authService) LoginUser(ctx context.Context, login, password string) (token string, err error) {
	user, err := s.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	err = util.CompareHashAndPassword(password, user.Password)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err = s.CreateToken(ctx, *user)
	if err != nil {
		return "", domain.ErrTokenCreation
	}

	return
}

func (s *authService) CreateToken(ctx context.Context, user domain.User) (token string, err error) {
	tokenPayload := domain.TokenPayload{
		UserID: user.ID,
	}

	token, err = s.tokenAdapter.CreateToken(tokenPayload)
	if err != nil {
		return "", domain.ErrTokenCreation
	}

	return
}
