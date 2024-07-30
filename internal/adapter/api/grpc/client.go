package grpc

import (
	"context"
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/api/grpc/interceptor"
	pb "github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/proto"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/port"
)

type ApiAdapter struct {
	address string
	conn    *grpc.ClientConn
	client  pb.VaultClient
}

// New creates a grpc client instance
func New(address string) (port.ApiAdapter, error) {
	return &ApiAdapter{
		address: address,
	}, nil
}

func (api *ApiAdapter) Open() error {
	conn, err := grpc.NewClient(api.address, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		ClientAuth:         tls.NoClientCert,
		InsecureSkipVerify: true,
	})), grpc.WithUnaryInterceptor(interceptor.AuthUnaryClientInterceptor))

	if err != nil {
		err = fmt.Errorf("failed to create grpc connection: %w", err)
		return err
	}
	api.client = pb.NewVaultClient(conn)
	return nil
}

func (api *ApiAdapter) Close() error {
	if api.conn != nil {
		return api.conn.Close()
	}
	return nil
}

func (api *ApiAdapter) SetToken(ctx context.Context, token string) (context.Context, error) {
	return metadata.AppendToOutgoingContext(ctx, "X-Auth-Token", token), nil
}
