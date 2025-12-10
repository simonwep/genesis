package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Health godoc
// @Summary      Health check
// @Description  Check if the API is healthy and responding
// @Tags         health
// @Produce      json
// @Success      200 "API is healthy"
// @Router       /health [get]
func Health(c *gin.Context) {

	// We assume, if the api is able to respond to this request, it is healthy.
	c.Status(http.StatusOK)
}
