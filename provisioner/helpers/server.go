package helpers

import (
	"log"
	"net"
	"time"
)

const (
    serverUpTimeout = 15 * time.Minute
    serverPingInterval = 10 * time.Second
)

func OnMcServerUp(ip string, serverUpFunc func()) {
	logger := log.Default()
	logger.Println("SERVER STARTING, IP: " + ip)
	timeoutTime := time.Now().Add(serverUpTimeout)
	// Loop to check server status
	for {
		if time.Now().After(timeoutTime) {
			logger.Printf("timeout: mc server starting up %s took too long\n", ip)
			return
		}

		if IsMCServerUp(ip, "25565") {
			logger.Printf("mc server %s is up! \n", ip)
			serverUpFunc()
			return
		} else {
			logger.Printf("server %s is down, pinging again in 10secs \n", ip)
		}

		time.Sleep(serverPingInterval) // Check every 10 seconds
	}
}

func IsMCServerUp(ip, port string) bool {
	address := net.JoinHostPort(ip, port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second) // 5-second timeout
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
