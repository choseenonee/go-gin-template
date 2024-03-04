package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"template/internal/delivery/middleware"
	"template/pkg/auth"
	"template/pkg/database/cached"
	"template/pkg/log"
)

func InitRouting(r *gin.Engine, db *sqlx.DB, logger *log.Logs, middlewareStruct middleware.Middleware, jwtUtils auth.JWTUtil, session cached.Session) {
	_ = RegisterPublicRouter(r, db, session, jwtUtils, logger)

	userRouter := r.Group("/user")
	userRouter.Use(middlewareStruct.Authorization())

	RegisterUserRouter(userRouter, db, session, jwtUtils, logger)
}
