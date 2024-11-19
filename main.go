package main

import (
	"CloudBook/internal/repository/dao"
	"CloudBook/internal/web"
	"CloudBook/internal/web/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:13306)/cloudbook"))
	if err != nil {
		// 只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	server := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))
	server.Use(middleware.NewLoginMiddleWareBuilder().Build())

	web.RegisterRoutes(server, db)

	err = server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
