package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/alancuriel/game-hosting-sass/provisioner/models"
	"github.com/alancuriel/game-hosting-sass/provisioner/services"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := log.Default()
	mcService, err := services.NewMinecraftLinodeProvisionService()
	if err != nil {
		logger.Fatal(err.Error())
	}

	r := gin.Default()

	r.Use(gin.Recovery())

	adminPswd := os.Getenv("PROVISIONER_ADMIN_PASS")
	if adminPswd == "" {
		logger.Fatal("could not start, No admin password found")
	}

	r.Use(gin.BasicAuth(gin.Accounts{
		"admin": adminPswd,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/v1/provision/mc", func(c *gin.Context) {
		var req *models.ProvisionMcServerRequest
		if err := c.BindJSON(&req); err != nil {
			return
		}

		if req.Instance == models.MINECRAFT_INSTANCE_INVALID ||
			req.Region == models.INVALID || req.Username == "" ||
			req.Owner == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		ip, err := mcService.Provision(req)
		if err != nil {
			logger.Println("error provisioning ", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.String(http.StatusCreated, ip)
	})

	r.GET("/v1/servers/mc/:owner", func(c *gin.Context) {
		owner := c.Param("owner")

		if strings.TrimSpace(owner) == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		servers, err := mcService.ListServersByOwner(owner)

		if err != nil {
			logger.Println("error listing servers ", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, servers)
	})

	r.DELETE("/v1/servers/mc/:id", func(c *gin.Context) {
		id := c.Param("id")

		if strings.TrimSpace(id) == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		err := mcService.DeleteServer(id)

		if err != nil {
			logger.Println("error deleting server ", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusNoContent)
	})

	r.POST("/v1/servers/mc/:id/announce", func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Message string `json:"message" binding:"required"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		if err := mcService.AnnounceMessage(id, req.Message); err != nil {
			logger.Printf("Failed to announce message: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	})

	r.Run()
}
