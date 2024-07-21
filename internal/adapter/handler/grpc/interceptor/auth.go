package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/port"
)

func AuthUnaryServerInterceptor(tokenAdapter port.TokenAdapter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// Skip interceptor for the following methods
		switch info.FullMethod {
		case "/service.Vault/LoginUser",
			"/service.Vault/RegisterUser":
			return handler(ctx, req)
		}

		var token string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			var bearer string
			values := md.Get("authorization")
			if len(values) > 0 {
				bearer = values[0]
			}
			if len(bearer) > 7 {
				authType := strings.ToLower(bearer[:6])
				if authType != "bearer" {
					return nil, status.Error(codes.Unauthenticated, "Unsupported authorization method")
				}
				token = bearer[7:]
			}
		}
		if len(token) == 0 {
			return nil, status.Error(codes.Unauthenticated, "Missing bearer token")
		}

		payload, err := tokenAdapter.VerifyToken(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "Invalid bearer token")
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(map[string]string{})
		}
		md.Append("X-User-Id", payload.UserID.String())
		ctx = metadata.NewIncomingContext(ctx, md)

		return handler(ctx, req)
	}
}
