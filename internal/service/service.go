package service

import (
	"context"
	"template/internal/model/entities"
)

type Public interface {
	CreateUser(ctx context.Context, userCreate entities.UserCreate) (string, string, error)
	LoginUser(ctx context.Context, userLogin entities.UserCreate) (string, string, error)
	Refresh(ctx context.Context, sessionID string) (string, string, error)
}

type User interface {
	GetMe(ctx context.Context, userID int) (entities.User, error)
	Delete(ctx context.Context, userID int, sessionID string) error
}
