package main

import (
	"CloudBook/internal/repository/dao"
	"CloudBook/internal/web"
	"CloudBook/internal/web/middleware"
	"CloudBook/pkg/ginx/middlewares/ratelimit"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
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

	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	server.Use(ratelimit.NewBuilderWithLimiter(redisClient, time.Second, 100).
		Build())

	//store := memstore.NewStore([]byte("hxzuKSQBmze7nJ6jssZJQWEPJJJ2trsxfD3nGpnzBuyXCdd6TS7ATS3SEAGWKzwd"),
	//	[]byte("3cAraCAc7BZxhpbFXDnQ4PuFezCUXhwDvBPKyhQH3HzH5pTmv4wGRzUUP2AmyRUD"))
	//store, err := redis.NewStore(10, "tcp", "localhost:6379", "",
	//	[]byte("hxzuKSQBmze7nJ6jssZJQWEPJJJ2trsxfD3nGpnzBuyXCdd6TS7ATS3SEAGWKzwd"),
	//	[]byte("3cAraCAc7BZxhpbFXDnQ4PuFezCUXhwDvBPKyhQH3HzH5pTmv4wGRzUUP2AmyRUD"))
	//if err != nil {
	//	panic(err)
	//}
	store := memstore.NewStore(
		[]byte("hxzuKSQBmze7nJ6jssZJQWEPJJJ2trsxfD3nGpnzBuyXCdd6TS7ATS3SEAGWKzwd"),
		[]byte("3cAraCAc7BZxhpbFXDnQ4PuFezCUXhwDvBPKyhQH3HzH5pTmv4wGRzUUP2AmyRUD"))

	server.Use(sessions.Sessions("mysession", store))

	//server.Use(middleware.NewLoginMiddleWareBuilder().
	//	IgnorePaths("/users/login").
	//	IgnorePaths("/users/signup").
	//	Build())
	server.Use(middleware.NewLoginJWTMiddleWareBuilder().
		IgnorePaths("/users/login").
		IgnorePaths("/users/signup").
		Build())

	web.RegisterRoutes(server, db)

	err = server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
