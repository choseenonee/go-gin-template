package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"template/internal/model/entities"
	"template/internal/service"
	"template/pkg/customerr"
)

type PublicHandler struct {
	service service.Public
}

func InitPublicHandler(service service.Public) PublicHandler {
	return PublicHandler{
		service: service,
	}
}

// @Summary Create user
// @Tags public
// @Accept  json
// @Produce  json
// @Param data body entities.UserCreate true "User create"
// @Success 200 {object} int "Successfully created user, returning JWT and Session"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /public/create [post]
func (p PublicHandler) CreateUser(c *gin.Context) {
	var userCreate entities.UserCreate

	if err := c.ShouldBindJSON(&userCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	userToken, sessionID, err := p.service.CreateUser(ctx, userCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"JWT": userToken, "Session": sessionID})
}

// @Summary Login in mobile user
// @Tags public
// @Accept  json
// @Produce  json
// @Param data body entities.UserCreate true "Mobile user login"
// @Success 200 {object} map[string]string "Successfully loginned user, returning JWT and Session"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /public/login [post]
func (p PublicHandler) LoginUser(c *gin.Context) {
	var userCreate entities.UserCreate

	if err := c.ShouldBindJSON(&userCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	userToken, sessionID, err := p.service.LoginUser(ctx, userCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"JWT": userToken, "Session": sessionID})
}

// @Summary Refresh tokens
// @Tags public
// @Accept  json
// @Produce  json
// @Param Session header string true "Session ID"
// @Success 200 {object} map[string]string "Successfully authorized, returning JWT and new session_id"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /public/refresh [post]
func (p PublicHandler) Refresh(c *gin.Context) {
	ctx := c.Request.Context()

	sessionID := c.GetHeader("Session")

	userToken, newSessionID, err := p.service.Refresh(ctx, sessionID)
	if err != nil {
		if errors.Is(err, customerr.UserNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": err.Error()})
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"JWT": userToken, "Session": newSessionID})
}
