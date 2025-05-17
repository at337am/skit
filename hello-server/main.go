package main

import (
	"hello-server/config"
	"hello-server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	
	r := gin.Default()
	
	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/*")

	routes.RegisterRoutes(r)

	r.Run(config.GetServerPort())
}
