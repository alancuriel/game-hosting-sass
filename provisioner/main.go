package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alancuriel/game-hosting-sass/provisioner/models"
	"github.com/alancuriel/game-hosting-sass/provisioner/services"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Default()
	mcService, err := services.NewMinecraftLinodeProvisionService()
	if err != nil {
		log.Fatal(err.Error())
	}

	r := gin.Default()

	r.Use(gin.Recovery())

	adminPswd := os.Getenv("PROVISIONER_ADMIN_PASS")
	if adminPswd == "" {
		log.Fatal("could not start, No admin password found")
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

		c.String(http.StatusCreated, ip)
	})

	r.Run()
}
