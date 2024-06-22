package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dimitargrozev5/expenses-go-1/internal/jwtutil"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/internal/sysinfo"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func registerDBNode() {

	// Open connection to DB Controller
	var opts = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(app.ControllerAddress, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create gRPC client
	client := models.NewDatabaseClient(conn)

	// Get system info
	props := sysinfo.Overview()

	// Create jwt
	jwt, err := jwtutil.Repo.Generate(jwt.MapClaims{})
	if err != nil {
		log.Fatal(err)
	}

	// Create context with metadata
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwt))
	ctxWithMeta := metadata.NewOutgoingContext(context.Background(), md)

	// Register node
	_, err = client.RegisterNode(ctxWithMeta, &props)
	if err != nil {
		log.Fatal(err)
	}
}
