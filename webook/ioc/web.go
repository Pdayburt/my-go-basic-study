package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"my-go-basic-study/webook/internal/web"
	"my-go-basic-study/webook/internal/web/middleware"
	"strings"
	"time"
)

func InitGin(userHandler *web.UserHandler, articleHandler *web.ArticleHandler, handlers []gin.HandlerFunc) *gin.Engine {
	server := gin.Default()
	userHandler.RegisterRouters(server)
	articleHandler.RegisterRouters(server)
	server.Use(handlers...)
	return server
}

func InitMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		InitCORS(),
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/users/login").IgnorePaths("/users/signup").
			IgnorePaths("/users/loginJWT").Build(),
	}
}

func InitCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
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
	})

}
