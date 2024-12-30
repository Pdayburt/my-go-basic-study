package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"my-go-basic-study/webook/internal/repository"
	"my-go-basic-study/webook/internal/repository/dao"
	"my-go-basic-study/webook/internal/service"
	"my-go-basic-study/webook/internal/web"
	"my-go-basic-study/webook/internal/web/middleware"
	"my-go-basic-study/webook/pkg/gin/middlewares/ratelimit"
	"net/http"
	"strings"
	"time"
)

func main() {

	/*server := initWebEngine()
	db := initDb()
	initUser(db, server)*/
	server := gin.Default()
	server.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello Gin")
		return
	})
	server.Run(":8080")

}

func initWebEngine() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods:     []string{"POST", "PUT", "PATCH", "GET"},
		AllowHeaders: []string{"Authorization", "content-type"},
		//ExposeHeaders 才能把header中的数据暴露出来，前端才能拿到x-jwt-token的值
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "http://localhost") {
				return true
			}
			return false
		},
		MaxAge: 12 * time.Hour,
	}))

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/login").IgnorePaths("/users/signup").Build())
	return server
}

func initUser(db *gorm.DB, server *gin.Engine) {
	userDao := dao.NewUserDao(db)
	userRepository := repository.NewUserRepository(userDao)
	userService := service.NewUserService(userRepository)
	userHandler := web.NewUserHandler(userService)
	userHandler.RegisterRouters(server)
}

func initDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		//panic 整个go routine 结束
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
