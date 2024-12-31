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
	healthRouter.UseRoutes(r)

	mcRouter := routes.NewMcRouter(mcService)
	mcRouter.UseRoutes(r)

	r.SetTrustedProxies([]string{})

	r.Run()
}
