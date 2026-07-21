package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sowmyavejerla13/url-shortener/internal/config"
	"github.com/sowmyavejerla13/url-shortener/internal/database"
	"github.com/sowmyavejerla13/url-shortener/internal/handler"
	"github.com/sowmyavejerla13/url-shortener/internal/repository"
	"github.com/sowmyavejerla13/url-shortener/internal/routes"
	"github.com/sowmyavejerla13/url-shortener/internal/service"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.NewPostgres(cfg)
	if err!=nil{
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo,cfg)
	userAuthHandler := handler.NewAuthHandler(userService)
	urlRepo := repository.NewURLRepository(db)
	urlService := service.NewURLService(urlRepo)
	urlHandler := handler.NewURLHandler(urlService,cfg.AppEnv)

	router := gin.Default()
	routes.SetupRoutes(router,userAuthHandler,urlHandler,cfg)

	log.Printf("Starting %s on port %s..\n",cfg.AppName,cfg.AppPort)

	if err:= router.Run(":"+cfg.AppPort);err!=nil{
		log.Fatal(err)
	}
}