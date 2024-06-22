package sysinfo

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/pbnjay/memory"
	"github.com/ricochet2200/go-disk-usage/du"
	"github.com/shirou/gopsutil/v3/cpu"
)

type SysinfoData struct {
	app    config.Config
	port   int
	nodeId int64
}

var data SysinfoData

func NewSysinfo(a config.Config, port int, nodeId int64) {
	data.app = a
	data.port = port
	data.nodeId = nodeId
}

var MB float64 = 1024 * 1024

func Overview() models.DBNodeData {

	// Get current address
	address := getLocalIP()
	if !data.app.GetInProduction() {
		address = "127.0.0.1"
	}

	// Add port
	address = fmt.Sprintf("%s:%d", address, data.port)

	// Get disk usage
	usage := du.NewDiskUsage(".")

	// Get CPU load percentage over 1 second interval
	percent, _ := cpu.Percent(time.Second, false)

	// Create props
	return models.DBNodeData{
		ID:             data.nodeId,
		Address:        address,
		TotalMemoryMB:  float64(memory.TotalMemory()) / MB,
		FreeMemoryMB:   float64(memory.FreeMemory()) / MB,
		TotalStorageMB: float64(usage.Size()) / MB,
		FreeStorageMB:  float64(usage.Available()) / MB,
		CpuLoadPercent: percent[0],
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
