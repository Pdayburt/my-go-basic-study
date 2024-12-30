package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"my-go-basic-study/webook/internal/web"
	"net/http"
	"strings"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (b *LoginJWTMiddlewareBuilder) IgnorePaths(paths string) *LoginJWTMiddlewareBuilder {
	b.paths = append(b.paths, paths)
	return b
}

func (b *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		//不需要校验
		for _, path := range b.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			//说明未登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var userClaim web.UserClaims
		token, err := jwt.ParseWithClaims(segs[1], &userClaim, func(*jwt.Token) (interface{}, error) {
			return []byte("a59b1c734e8f2d5a967b3c841e5f9a2d6b4"), nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("claims", userClaim)
	}

}
