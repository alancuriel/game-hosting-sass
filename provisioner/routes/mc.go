package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/alancuriel/game-hosting-sass/provisioner/models"
	"github.com/alancuriel/game-hosting-sass/provisioner/services"
	"github.com/gin-gonic/gin"
)

type McRouter struct {
	mcProvisionerService services.MinecraftProvisionService
	logger               *log.Logger
}

func NewMcRouter(service services.MinecraftProvisionService) *McRouter {
	return &McRouter{
		mcProvisionerService: service,
		logger:               log.Default(),
	}
}

func (r *McRouter) Provision(c *gin.Context) {
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

	ip, err := r.mcProvisionerService.Provision(req)
	if err != nil {
		r.logger.Println("error provisioning ", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusCreated, ip)
}

func (r *McRouter) ListByOwner(c *gin.Context) {
	owner := c.Param("owner")

	if strings.TrimSpace(owner) == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	servers, err := r.mcProvisionerService.ListServersByOwner(owner)

	if err != nil {
		r.logger.Println("error listing servers ", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, servers)
}

func (r *McRouter) DeleteServer(c *gin.Context) {
	id := c.Param("id")

	if strings.TrimSpace(id) == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := r.mcProvisionerService.DeleteServer(id)

	if err != nil {
		r.logger.Println("error deleting server ", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (r *McRouter) Announce(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := r.mcProvisionerService.AnnounceMessage(id, req.Message); err != nil {
		r.logger.Printf("Failed to announce message: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
