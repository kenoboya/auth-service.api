package service

import (
	"auth-service/internal/model"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Users interface {
	Create(ctx context.Context, user model.UserSignUp) (bson.ObjectID, error)
	GetByLogin(ctx context.Context, login string) (model.User, error)
}
