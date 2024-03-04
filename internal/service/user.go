package service

import (
	"context"
	"github.com/spf13/viper"
	"template/internal/model/entities"
	"template/internal/repository"
	"template/pkg/auth"
	"template/pkg/config"
	"template/pkg/database/cached"
	"template/pkg/log"
	"time"
)

type userService struct {
	userRepo        repository.User
	timeoutDuration time.Duration
	session         cached.Session
	jwt             auth.JWTUtil
	logger          *log.Logs
}

func InitUserService(
	userRepo repository.User,
	session cached.Session,
	jwt auth.JWTUtil,
	logger *log.Logs,
) User {
	return userService{
		userRepo:        userRepo,
		timeoutDuration: time.Duration(viper.GetInt(config.TimeOut)) * time.Millisecond,
		session:         session,
		jwt:             jwt,
		logger:          logger,
	}
}

func (u userService) GetMe(ctx context.Context, userID int) (entities.User, error) {
	return u.userRepo.Get(ctx, userID)
}

func (u userService) Delete(ctx context.Context, userID int, sessionID string) error {
	err := u.session.Delete(ctx, userID, sessionID)
	if err != nil {
		return err
	}
	return u.userRepo.Delete(ctx, userID)
}
