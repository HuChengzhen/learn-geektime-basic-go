package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"learn-geektime-basic-go/webook/internal/web"
	"net/http"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(context *gin.Context) {
		for _, path := range l.paths {
			if context.Request.URL.Path == path {
				return
			}
		}

		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[7:]

		claims := &web.UserClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("zonUXJUU5%XmP6wkH^X%W7l%sNM0dPvI"), nil
		})

		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid || claims.Uid == 0 {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		context.Set("claims", claims)

		fmt.Println(token)
	}
}
