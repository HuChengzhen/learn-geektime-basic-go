package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"learn-geektime-basic-go/webook/internal/web"
)

func main() {
	server := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"authorization", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}
	server.Use(cors.New(config))

	u := web.NewUserHandler()
	u.RegisterRoutes(server)
	server.Run(":8080")
}
