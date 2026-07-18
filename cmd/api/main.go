package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sowmyavejerla13/url-shortener/internal/config"
	"github.com/sowmyavejerla13/url-shortener/internal/routes"
)

func main() {
	cfg := config.LoadConfig()
	router := gin.Default()
	routes.RegisterRoutes(router)
	log.Printf("Starting %s on port %s..\n",cfg.AppName,cfg.AppPort)

	if err:= router.Run(":"+cfg.AppPort);err!=nil{
		log.Fatal(err)
	}
}