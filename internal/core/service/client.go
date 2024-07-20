package service

import (
	"context"
	
	"github.com/fishus/go-advanced-gophkeeper/internal/core/port"
)

type clientService struct {
	apiAdapter port.ApiAdapter
	token      string
}

func NewClientService(
	apiAdapter port.ApiAdapter,
) port.ClientService {
	return &clientService{
		apiAdapter: apiAdapter,
	}
}

func (s *clientService) Setup(ctx context.Context) error {
	return s.apiAdapter.Open()
}

func (s *clientService) Teardown(ctx context.Context) error {
	return s.apiAdapter.Close()
}

func (s *clientService) SetToken(token string) port.ClientService {
	s.token = token
	return s
}

func (s *clientService) UserLogin(ctx context.Context, login, password string) (token string, err error) {
	token, err = s.apiAdapter.LoginUser(ctx, login, password)
	if err == nil {
		s.token = token
	}
	return
}

func (s *clientService) UserRegister(ctx context.Context, login, password string) (token string, err error) {
	token, err = s.apiAdapter.RegisterUser(ctx, login, password)
	if err == nil {
		s.token = token
	}
	return
}
