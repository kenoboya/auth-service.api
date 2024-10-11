package repo

import (
	"auth-service/internal/model"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repositories struct {
	Users Users
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Users: NewUsersRepo(db),
	}
}

type Users interface {
	Create(ctx context.Context, user model.User) (bson.ObjectID, error)
	GetByLogin(ctx context.Context, login string) (model.User, error)
	GetSessionInfoByUserID(ctx context.Context, userID bson.ObjectID) (model.UserResponse, error)
}
