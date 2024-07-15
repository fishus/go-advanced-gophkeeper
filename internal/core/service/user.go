package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/port"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/util"
)

type userService struct {
	userRepo port.UserRepository
}

func NewUserService(userRepo port.UserRepository) port.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) RegisterUser(ctx context.Context, user domain.User) (*domain.User, error) {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("user not created: %w", err)
	}

	user.Password = hashedPassword

	usr, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("user not created: %w", err)
	}

	return usr, nil
}
