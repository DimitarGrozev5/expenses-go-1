package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/pbnjay/memory"
	"github.com/ricochet2200/go-disk-usage/du"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	// Get current address
	address := getLocalIP()
	if !app.InProduction {
		address = "127.0.0.1"
	}

	// Add port
	address = fmt.Sprintf("%s:%d", address, *port)

	// Get disk usage
	usage := du.NewDiskUsage(".")

	// Create props
	props := models.DBNodeData{
		Address:      address,
		TotalMemory:  float64(memory.TotalMemory()),
		FreeMemory:   float64(memory.FreeMemory()),
		TotalStorage: float64(usage.Size()),
		FreeStorage:  float64(usage.Available()),
	}

	_, err = client.RegisterNode(context.Background(), &props)
	if err != nil {
		log.Fatal(err)
	}
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP.String()
}
