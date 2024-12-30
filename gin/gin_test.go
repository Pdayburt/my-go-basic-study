package gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

func TestGin(t *testing.T) {

	server := gin.Default()

	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello gin")
	})

	server.GET("/user/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		ctx.String(http.StatusOK, "hello 这是参数路由 :"+name)
	})

	server.GET("/order", func(ctx *gin.Context) {
		value := ctx.Query("id")
		ctx.String(http.StatusOK, "query id = "+value)
	})

	server.GET("/view/*.html", func(ctx *gin.Context) {
		page := ctx.Param(".html")
		ctx.String(http.StatusOK, "hello 这是通配符路由 :"+page)

	})

	server.POST("/post", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "post gin")
	})

	server.Run(":8080")

}
