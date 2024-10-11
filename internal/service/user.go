package service

import (
	"auth-service/internal/model"
	"auth-service/internal/repository/cache"
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
	repo            repo.Users
	hasher          hash.PasswordHasher
	tokenManager    auth.TokenManager
	cache           cache.SessionCache
	sessionTokenTTL time.Duration
}

func NewUsersService(repo repo.Users, hasher hash.PasswordHasher,
	tokenManager auth.TokenManager,
	cache cache.SessionCache,
	sessionTokenTTL time.Duration) *UsersService {
	return &UsersService{
		repo:            repo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		cache:           cache,
		sessionTokenTTL: sessionTokenTTL,
	}
}

func (s *UsersService) SignUp(ctx context.Context, user model.UserSignUp) (model.Session, error) {
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

	return s.createSession(ctx, id.String())
}
func (s *UsersService) SignIn(ctx context.Context, requestSignIn model.UserSignIn) (model.Session, error) {
	user, err := s.repo.GetByLogin(ctx, requestSignIn.Login)
	if err != nil {
		return model.Session{}, err
	}

	if err := s.hasher.Compare(user.Password, requestSignIn.Password); err != nil {
		return model.Session{}, err
	}

	return s.createSession(ctx, user.UserID.String())
}
func (s *UsersService) Verify(ctx context.Context, sessionToken string) (model.UserResponse, error) {
	userID, err := s.cache.GetSession(ctx, sessionToken)
	if err != nil {
		return model.UserResponse{}, err
	}
	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return model.UserResponse{}, model.ErrInvalidFormatForConvertObjectID
	}
	userResponse, err := s.repo.GetSessionInfoByUserID(ctx, objectID)
	if err != nil {
		return model.UserResponse{}, err
	}
	return userResponse, nil
}

func (s *UsersService) createSession(ctx context.Context, id string) (model.Session, error) {
	sessionToken, err := s.tokenManager.GenerateSessionToken(s.hasher)
	if err != nil {
		logger.Error(
			zap.String("action", "GenerateSessionToken()"),
			zap.Error(err),
		)
		return model.Session{}, err
	}

	session := model.Session{
		Token:     sessionToken,
		ExpiresAt: time.Now().Add(s.sessionTokenTTL),
	}

	if err := s.cache.SetSession(ctx, session, id); err != nil {
		return model.Session{}, err
	}

	return session, nil
}
