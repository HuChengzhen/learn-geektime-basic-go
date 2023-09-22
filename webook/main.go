package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"learn-geektime-basic-go/webook/internal/repository"
	"learn-geektime-basic-go/webook/internal/repository/dao"
	"learn-geektime-basic-go/webook/internal/service"
	"learn-geektime-basic-go/webook/internal/web"
)

func main() {
	server := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"authorization", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}
	server.Use(cors.New(config))

	db, err := gorm.Open(mysql.Open("root:root@tcp(192.168.88.131:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	err = dao.InitTable(db)
	if err != nil {
		panic("failed to init table")
	}
	userDAO := dao.NewUserDAO(db)
	rp := repository.NewUserRepository(userDAO)
	svc := service.NewUserService(rp)
	u := web.NewUserHandler(svc)
	u.RegisterRoutes(server)
	server.Run(":8080")
}
