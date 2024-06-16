package rpcserver

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

func (s DatabaseServer) AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// Create context
	userCtx := ctx

	// Skip auth for some methods
	if !(strings.HasSuffix(info.FullMethod, "/Authenticate") || strings.HasSuffix(info.FullMethod, "/RegisterNode")) {
		// authentication (token verification)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errMissingMetadata
		}
		auth := md["authorization"]

		if len(auth) < 1 {
			return nil, errInvalidToken
		}
		token := strings.TrimPrefix(auth[0], "Bearer ")

		// Parse Token
		t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return s.App.JWTSecretKey, nil
		})
		if err != nil {
			return nil, errInvalidToken
		}

		// Get claims
		claims, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errInvalidToken
		}

		// Store token details
		userCtx = context.WithValue(ctx, "userKey", claims["userKey"])
		userCtx = context.WithValue(userCtx, "dbVersion", claims["dbVersion"])
	}

	m, err := handler(userCtx, req)
	if err != nil {
		s.App.ErrorLog.Fatalf("RPC failed with error: %v", err)
	}
	return m, err
}
