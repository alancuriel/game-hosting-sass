package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthRouter struct {
}

func NewHealthRouter() *HealthRouter {
	return &HealthRouter{}
}

func (h *HealthRouter) UseRoutes(e *gin.Engine) {
    e.GET("/ping", h.healthCheck)
}

func (h *HealthRouter) healthCheck(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
