package helpers

import (
	"log"
	"net"
	"time"
)

func PingMcServer(ip string) {
	logger := log.Default()
	logger.Println("SERVER STARTING, IP: " + ip)

	// Loop to check server status
	for {
		if IsServerUp(ip, "25565") {
			logger.Println("Server is UP!")
			return
		} else {
			logger.Println("Server is DOWN!")
		}

		time.Sleep(10 * time.Second) // Check every 10 seconds
	}
}

func IsServerUp(ip, port string) bool {
	address := net.JoinHostPort(ip, port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second) // 5-second timeout
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
