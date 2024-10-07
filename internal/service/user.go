package service

import (
	"auth-service/internal/model"
	repo "auth-service/internal/repository/mongo"
	"auth-service/pkg/auth"
	"auth-service/pkg/hash"
	logger "auth-service/pkg/logger/zap"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"
)

type UsersService struct {
	repo         repo.Users
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	cache
}

func NewUsersService(repo repo.Users, hasher hash.PasswordHasher, tokenManager auth.TokenManager) *UsersService {
	return &UsersService{
		repo:         repo,
		hasher:       hasher,
		tokenManager: tokenManager,
	}
}

func (s *UsersService) Create(ctx context.Context, user model.UserSignUp) (model.Session, error) {
	id, err := s.repo.Create(ctx, model.User{
		UserID:       bson.NewObjectID(),
		Username:     user.Username,
		Email:        user.Email,
		Password:     user.Password,
		Role:         model.RoleCustomer,
		RegisteredAt: time.Now(),
		Person: model.Person{
			FirstName: user.Person.FirstName,
			LastName:  user.Person.LastName,
			Age:       user.Person.Age,
		},
	})
	if err != nil {
		return model.Session{}, err
	}
	session, err := s.tokenManager.GenerateSessionToken(s.hasher)
	if err != nil {
		logger.Error(
			zap.String("action", "GenerateSessionToken()"),
			zap.Error(err),
		)
		return model.Session{}, err
	}
	
}
func (s *UsersService) GetByLogin(ctx context.Context, login string) (model.User, error)
