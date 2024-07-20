package port

import "context"

type ApiAdapter interface {
	Open() error
	Close() error
	LoginUser(ctx context.Context, login, password string) (token string, err error)
	RegisterUser(ctx context.Context, login, password string) (token string, err error)
}
