package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// unaryInterceptor is an example unary interceptor.
func authInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Check for auth token in context
	token := app.Session.GetString(ctx, "user_token")

	// Define variable
	ctxWithMeta := ctx

	// If token is set
	if len(token) > 0 {

		// Add auth token to context
		md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", token))
		ctxWithMeta = metadata.NewOutgoingContext(ctx, md)
	}

	err := invoker(ctxWithMeta, method, req, reply, cc, opts...)
	return err
}
