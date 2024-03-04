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

func RegisterUserRouter(userRouter *gin.RouterGroup, db *sqlx.DB, session cached.Session, jwt auth.JWTUtil, logger *log.Logs) {
	userRepo := user.InitUserRepo(db)

	userService := service.InitUserService(userRepo, session, jwt, logger)
	userHandler := handlers.InitUserHandler(userService, session)

	userRouter.GET("/me", userHandler.GetMe)
	userRouter.GET("/delete", userHandler.Delete)
}
