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

func (h *HealthRouter) HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
