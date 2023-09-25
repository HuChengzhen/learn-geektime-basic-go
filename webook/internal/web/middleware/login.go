package middleware

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(context *gin.Context) {
		for _, path := range l.paths {
			if context.Request.URL.Path == path {
				return
			}
		}

		sess := sessions.Default(context)
		if sess == nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		id := sess.Get("userId")
		if id == nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := sess.Get("update_time")
		now := time.Now()

		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}

		updateTimeVal, ok := updateTime.(time.Time)

		if !ok {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if now.Sub(updateTimeVal) >= time.Minute {
			sess.Set("update_time", now)
			sess.Save()
			return
		}

		fmt.Println(id)
	}
}
