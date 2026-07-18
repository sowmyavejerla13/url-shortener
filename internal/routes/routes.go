package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sowmyavejerla13/url-shortener/internal/handler"
)

func RegisterRoutes(router *gin.Engine){
	router.GET("/health", handler.Health)
}