package middleware

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddleWareBuilder struct {
}

func NewLoginMiddleWareBuilder() *LoginMiddleWareBuilder {
	return &LoginMiddleWareBuilder{}
}

func (l *LoginMiddleWareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		fmt.Println("Build-----------------")
		if ctx.Request.URL.Path == "/users/login" ||
			ctx.Request.URL.Path == "/users/signup" {
			return
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		println(id)
		if id == nil {
			//没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
