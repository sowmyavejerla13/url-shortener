package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sowmyavejerla13/url-shortener/internal/config"
	"github.com/sowmyavejerla13/url-shortener/internal/handler"
	"github.com/sowmyavejerla13/url-shortener/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(
	router *gin.Engine,
	authHandler *handler.AuthHandler,
	urlHandler *handler.URLHandler,
	cfg *config.Config,
) {
	router.GET("/:shortCode", urlHandler.Redirect)
	router.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api/v1")
	router.GET("/health", handler.Health)

	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.GET("/me", authHandler.Me)
		protected.GET("/urls", urlHandler.GetUserURLs)
		protected.POST("/shorten", urlHandler.Create)
		protected.DELETE("/urls/:id", urlHandler.Delete)
	}

}
