package port

import "context"

type ClientService interface {
	Setup(context.Context) error
	Teardown(context.Context) error
	SetToken(token string) ClientService
	UserLogin(ctx context.Context, login, password string) (token string, err error)
	UserRegister(ctx context.Context, login, password string) (token string, err error)
}
