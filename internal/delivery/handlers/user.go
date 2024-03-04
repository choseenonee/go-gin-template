package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"template/internal/delivery/middleware"
	"template/internal/service"
	"template/pkg/database/cached"
)

type UserHandler struct {
	service service.User
	session cached.Session
}

func InitUserHandler(service service.User, session cached.Session) UserHandler {
	return UserHandler{
		service: service,
		session: session,
	}
}

// @Summary Get user data
// @Tags user
// @Accept  json
// @Produce  json
// @Param Session header string true "Session ID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} []entities.User "Successfully response with user data"
// @Failure 400 {object} map[string]string "JWT is absent or invalid input"
// @Failure 403 {object} map[string]string "JWT is invalid or expired"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/me [get]
func (u UserHandler) GetMe(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetInt(middleware.CUserID)

	user, err := u.service.GetMe(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": err})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Delete user
// @Tags user
// @Accept  json
// @Produce  json
// @Param Session header string true "Session ID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} map[string]string "Successfully response"
// @Failure 400 {object} map[string]string "JWT is absent or invalid input"
// @Failure 403 {object} map[string]string "JWT is invalid or expired"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/delete [get]
func (u UserHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	sessionID := c.GetString(middleware.CSessionID)
	userID := c.GetInt(middleware.CUserID)

	err := u.service.Delete(ctx, userID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"detail": "successfully!"})
}
