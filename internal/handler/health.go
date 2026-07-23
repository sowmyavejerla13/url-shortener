package handler

import "github.com/gin-gonic/gin"

// Health godoc
//
// @Summary Health Check
// @Description Returns the health status of the application.
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func Health(c *gin.Context) {

	c.JSON(200, gin.H{
		"status":  "UP",
		"service": "url-shorterner",
		"version": "1.0.0",
	})
}
