package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"template/internal/delivery/handlers"
	"template/internal/repository/user"
	"template/internal/service"
	"template/pkg/auth"
	"template/pkg/database/cached"
	"template/pkg/log"
)

func RegisterPublicRouter(r *gin.Engine, db *sqlx.DB, session cached.Session, jwtUtil auth.JWTUtil, logger *log.Logs) *gin.RouterGroup {
	publicRouter := r.Group("/public")

	userRepo := user.InitUserRepo(db)

	publicService := service.InitPublicService(userRepo, session, jwtUtil, logger)
	publicHandler := handlers.InitPublicHandler(publicService)

	publicRouter.POST("/create", publicHandler.CreateUser)
	publicRouter.POST("/login", publicHandler.LoginUser)

	publicRouter.POST("/refresh", publicHandler.Refresh)

	return publicRouter
}
