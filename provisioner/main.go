package main

import (
	"log"
	"os"

	"github.com/alancuriel/game-hosting-sass/provisioner/routes"
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

	healthRouter := routes.NewHealthRouter()
	r.GET("/ping", healthRouter.HealthCheck)

	mcRouter := routes.NewMcRouter(mcService)

	r.POST("/v1/provision/mc", mcRouter.Provision)
	r.GET("/v1/servers/mc/:owner", mcRouter.ListByOwner)
	r.DELETE("/v1/servers/mc/:id", mcRouter.DeleteServer)
	r.POST("/v1/servers/mc/:id/announce", mcRouter.Announce)

	r.SetTrustedProxies([]string{})

	r.Run()
}
