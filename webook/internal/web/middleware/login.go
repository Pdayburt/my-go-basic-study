package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddleBuilder struct {
	paths []string
}

func NewLoginMiddleBuilder() *LoginMiddleBuilder {
	return &LoginMiddleBuilder{}
}

func (b *LoginMiddleBuilder) IgnorePaths(paths string) *LoginMiddleBuilder {
	b.paths = append(b.paths, paths)
	return b
}

func (b *LoginMiddleBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//不需要校验
		for _, path := range b.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		sess.Options(sessions.Options{
			MaxAge: 60 * 60,
		})
		updateTime := sess.Get("update_time")
		now := time.Now().UnixMilli()
		//update_time没有 说明是登录后第一次进入
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
		//update_time存在
		updateTimeVal, _ := updateTime.(int64)
		if now-updateTimeVal > 59*60*1000 {
			sess.Set("update_time", now)
			sess.Save()
			return
		}

	}
}
