package helpers

import (
	"log"
	"net"
	"time"
)

func PingMcServer(ip string, callback func()) {
	logger := log.Default()
	logger.Println("SERVER STARTING, IP: " + ip)
	timeoutTime := time.Now().Add(15 * time.Minute)
	// Loop to check server status
	for {
		if time.Now().After(timeoutTime) {
			logger.Printf("TIMEOUT mc server starting on %s took too long\n", ip)
			return
		}

		if IsServerUp(ip, "25565") {
			logger.Printf("Server %s is Up! \n", ip)
			callback()
			return
		} else {
			logger.Printf("Server %s is Down, pinging in 10secs \n", ip)
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
