package main

import (
	"encoding/gob"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"learn-geektime-basic-go/webook/internal/repository"
	"learn-geektime-basic-go/webook/internal/repository/dao"
	"learn-geektime-basic-go/webook/internal/service"
	"learn-geektime-basic-go/webook/internal/web"
	"learn-geektime-basic-go/webook/internal/web/middleware"
	"time"
)

func main() {
	db := initDB()
	server := initServer()
	initUser(db, server)
	server.Run(":8080")
}

func initUser(db *gorm.DB, server *gin.Engine) {
	userDAO := dao.NewUserDAO(db)
	rp := repository.NewUserRepository(userDAO)
	svc := service.NewUserService(rp)
	u := web.NewUserHandler(svc)
	u.RegisterRoutes(server)
}

func initServer() *gin.Engine {
	server := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"authorization", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}
	config.ExposeHeaders = []string{"x-jwt-token"}
	server.Use(cors.New(config))

	//store := cookie.NewStore([]byte("secret"))
	//store := memstore.NewStore([]byte("e%YTe2tmIlAFkO26a21dW!PtboY7kS*4"), []byte("uOI@Sum7m!%WHbscKrrml$^!2Ww#rURn"))

	store, err := redis.NewStore(10,
		"tcp",
		"192.168.88.131:6379",
		"",
		[]byte("e%YTe2tmIlAFkO26a21dW!PtboY7kS*4"),
		[]byte("uOI@Sum7m!%WHbscKrrml$^!2Ww#rURn"),
	)

	gob.Register(time.Now())

	store.Options(sessions.Options{
		// 10 minute
		MaxAge: 10 * 60,
	})

	if err != nil {
		panic("redis connect failed.")
	}

	server.Use(sessions.Sessions("mysession", store))

	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").
		Build(),
	)
	return server
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(192.168.88.131:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	err = dao.InitTable(db)
	if err != nil {
		panic("failed to init table")
	}
	return db
}
