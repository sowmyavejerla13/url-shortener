package handler

import "github.com/gin-gonic/gin"

func Health(c *gin.Context) {

	c.JSON(200, gin.H{
		"status" :"UP",
		"service":"url-shorterner",
		"version":"1.0.0",
	})
}