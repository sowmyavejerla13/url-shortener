package main

// @title URL Shortener API
// @version 1.0
// @description Production-ready URL Shortener API built with Go, Gin and PostgreSQL.
// @termsOfService http://swagger.io/terms/

// @contact.name Sowmya Vejerla
// @contact.url https://github.com/sowmyavejerla13
// @contact.email sowmya.vejerla@gmail.com

// @license.name MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/sowmyavejerla13/url-shortener/docs"

	"github.com/gin-gonic/gin"
	"github.com/sowmyavejerla13/url-shortener/internal/config"
	"github.com/sowmyavejerla13/url-shortener/internal/database"
	"github.com/sowmyavejerla13/url-shortener/internal/handler"
	"github.com/sowmyavejerla13/url-shortener/internal/repository"
	"github.com/sowmyavejerla13/url-shortener/internal/routes"
	"github.com/sowmyavejerla13/url-shortener/internal/service"
	"github.com/sowmyavejerla13/url-shortener/internal/utils"
	"github.com/sowmyavejerla13/url-shortener/pkg/logger"
)

func main() {
	log := logger.New()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error(
			"failed to load configuration",
			"error", err,
		)
		os.Exit(1)
	}

	log.Info("configuration loaded")

	// Connect database
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Error(
			"failed to connect to PostgreSQL",
			"error", err,
		)
		os.Exit(1)
	}

	log.Info("connected to PostgreSQL")
	defer db.Close()

	// Dependency Injection
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg)
	userHandler := handler.NewAuthHandler(userService)

	urlRepo := repository.NewURLRepository(db)
	urlService := service.NewURLService(urlRepo, utils.GenerateShortCode)
	urlHandler := handler.NewURLHandler(urlService)

	router := gin.Default()
	routes.SetupRoutes(router, userHandler, urlHandler, cfg)

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	// Start server
	go func() {
		log.Info(
			"starting server",
			"app", cfg.AppName,
			"port", cfg.AppPort,
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(
				"server failed",
				"error", err,
			)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-quit

	log.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error(
			"graceful shutdown failed",
			"error", err,
		)
	} else {
		log.Info("server stopped gracefully")
	}

	log.Info("closing database connection")
	db.Close()

	log.Info("application shutdown complete")
}
