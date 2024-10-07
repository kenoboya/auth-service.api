package service

import (
	"auth-service/internal/model"
	repo "auth-service/internal/repository/mongo"
	"auth-service/pkg/hash"
	"context"
)

type Services struct {
	Users Users
}

type Deps struct {
	repo         repo.Repositories
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	cache
}

func NewServices(deps Deps) *Services {
	return
}

type Users interface {
	Create(ctx context.Context, user model.UserSignUp) (model.Session, error)
	GetByLogin(ctx context.Context, login string) (model.User, error)
}
