package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthUnaryClientInterceptor(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Skip interceptor for the following methods
	switch method {
	case "/service.Vault/LoginUser",
		"/service.Vault/RegisterUser":
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			md.Delete("X-Auth-Token")
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	var token string
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		values := md.Get("X-Auth-Token")
		if len(values) > 0 {
			token = values[0]
		}
		md.Delete("X-Auth-Token")
	}

	return invoker(metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+token), method, req, reply, cc, opts...)
}
