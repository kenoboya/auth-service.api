package service

import (
	"auth-service/internal/config"
	"auth-service/internal/model"
	"auth-service/internal/repository/cache"
	repo "auth-service/internal/repository/mongo"
	"auth-service/pkg/auth"
	"auth-service/pkg/hash"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Deps struct {
	repo            *repo.Repositories
	hasher          hash.PasswordHasher
	tokenManager    auth.TokenManager
	cache           cache.SessionCache
	sessionTokenTTL time.Duration
}

func NewDeps(repo *repo.Repositories, cfg *config.Config, client *redis.Client) (*Deps, error) {
	hasher := hash.NewSHA256Hasher(cfg.Auth.Salt)
	tokenManager, err := auth.NewManager(cfg.Auth.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("tokenManager: %w", err)
	}
	cache := cache.NewRedisCache(client)
	return &Deps{
		repo:            repo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		cache:           cache,
		sessionTokenTTL: cfg.Auth.TokenTTL,
	}, nil
}

type Services struct {
	Users Users
}

func NewServices(deps *Deps) *Services {
	return &Services{
		Users: NewUsersService(deps.repo.Users,
			deps.hasher, deps.tokenManager, deps.cache,
			deps.sessionTokenTTL),
	}
}

type Users interface {
	SignUp(ctx context.Context, user model.UserSignUp) (model.Session, error)
	SignIn(ctx context.Context, requestSignIn model.UserSignIn) (model.Session, error)
}
