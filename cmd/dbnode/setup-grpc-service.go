package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/dimitargrozev5/expenses-go-1/internal/dbnoderpc"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"google.golang.org/grpc"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 3003, "The server port")
)

func setupGrpcService() {

	// Start listening on specified port
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Setup grpc server
	var opts []grpc.ServerOption
	// if *tls {
	// 	if *certFile == "" {
	// 		*certFile = data.Path("x509/server_cert.pem")
	// 	}
	// 	if *keyFile == "" {
	// 		*keyFile = data.Path("x509/server_key.pem")
	// 	}
	// 	creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
	// 	if err != nil {
	// 		log.Fatalf("Failed to generate credentials: %v", err)
	// 	}
	// 	opts = []grpc.ServerOption{grpc.Creds(creds)}
	// }

	// Create service
	databaseServer := dbnoderpc.NewService(&app)

	// Register server
	dbnoderpc.NewDatabaseServer(databaseServer)

	// Add JWT token interceptor
	opts = append(opts, grpc.UnaryInterceptor(dbnoderpc.Server.AuthInterceptor))

	// Create server
	grpcServer := grpc.NewServer(opts...)

	// Register server
	models.RegisterDatabaseServer(grpcServer, databaseServer)

	// Start grpc server
	fmt.Printf("Starting gRPC server on port %d", *port)
	grpcServer.Serve(lis)
}
