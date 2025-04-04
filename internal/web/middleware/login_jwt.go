package middleware

import (
	"CloudBook/internal/web"
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
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
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("3cAraCAc7BZxhpbFXDnQ4PuFezCUXhwDvBPKyhQH3HzH5pTmv4wGRzUUP2AmyRUD"), nil
		})
		//token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//	return []byte("3cAraCAc7BZxhpbFXDnQ4PuFezCUXhwDvBPKyhQH3HzH5pTmv4wGRzUUP2AmyRUD"), nil
		//})
		if err != nil {
			// 没登录,或者系统错误 Bearer xxx1234
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 解析时会自己验证过期时间
		// err 为 nil, token 肯定不为 nil, 这是约定俗成的
		if token == nil || !token.Valid || claims.Uid == 0 {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.UserAgent != ctx.Request.UserAgent() {
			// 严重的安全问题
			// 你是要监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 每十秒钟刷新一次
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("3cAraCAc7BZxhpbFXDnQ4PuFezCUXhwDvBPKyhQH3HzH5pTmv4wGRzUUP2AmyRUD"))
			if err != nil {
				// 记录日志
				log.Println("jwt 续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}

		ctx.Set("claims", claims)
	}
}
