package main

import (
	"github.com/gin-gonic/gin"
	"learn-geektime-basic-go/webook/internal/web"
)

func main() {
	server := gin.Default()
	u := &web.UserHandler{}
	u.RegisterRoutes(server)
	server.Run(":8080")
}
