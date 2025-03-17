package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

// LoginJWTMiddleWareBuilder 扩展性
type LoginJWTMiddleWareBuilder struct {
	paths []string
}

func NewLoginJWTMiddleWareBuilder() *LoginJWTMiddleWareBuilder {
	return &LoginJWTMiddleWareBuilder{}
}

func (l *LoginJWTMiddleWareBuilder) IgnorePaths(path string) *LoginJWTMiddleWareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddleWareBuilder) Build() gin.HandlerFunc {
	// 用 GO 的方式编码解码
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		// 现在用 JWT 校验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//segs := strings.SplitN(tokenHeader, " ", 2)
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			// 没登录,有人瞎搞
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("3cAraCAc7BZxhpbFXDnQ4PuFezCUXhwDvBPKyhQH3HzH5pTmv4wGRzUUP2AmyRUD"), nil
		})
		if err != nil {
			// 没登录,或者系统错误 Bearer xxx1234
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// err 为 nil, token 肯定不为 nil, 这是约定俗成的
		if token == nil || !token.Valid {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
