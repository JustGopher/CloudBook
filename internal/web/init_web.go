package web

import (
	"CloudBook/internal/repository"
	"CloudBook/internal/repository/dao"
	"CloudBook/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

func RegisterRoutes(server *gin.Engine, db *gorm.DB) *gin.Engine {
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	registerUsersRoutes(server, db)
	return server
}

func registerUsersRoutes(server *gin.Engine, db *gorm.DB) {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := NewUserHandler(svc)
	u.RegisterUserRoutes(server)
}
