package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/alancuriel/game-hosting-sass/provisioner/models"
	"github.com/alancuriel/game-hosting-sass/provisioner/services"
	"github.com/gin-gonic/gin"
)

func main() {
	mcService, err := services.NewMinecraftLinodeProvisionService()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/v1/provision/mc", func(c *gin.Context) {
		var req *models.ProvisionMcServerRequest
		if err := c.BindJSON(&req); err != nil {
			return
		}

		if req.Instance == models.MINECRAFT_INSTANCE_INVALID ||
			req.Region == models.INVALID || req.Username == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		ip, err := mcService.Provision(req.Instance, req.Region, req.Username)
		if err != nil {
			log.Default().Println("error provisioning ", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		go PingMcServer(ip)

		c.String(http.StatusCreated, ip)
	})

	r.Run()
}

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
