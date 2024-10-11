package repo

import (
	"auth-service/internal/model"
	"auth-service/pkg/database/mongodb"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UsersRepo struct {
	db *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
	return &UsersRepo{
		db: db.Collection(usersCollection),
	}
}

func (r *UsersRepo) Create(ctx context.Context, user model.User) (bson.ObjectID, error) {
	result, err := r.db.InsertOne(ctx, user)
	if err != nil {
		if mongodb.IsDuplicate(err) {
			return bson.ObjectID{}, model.ErrUserAlreadyExists
		}
		return bson.ObjectID{}, err
	}

	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return bson.ObjectID{}, model.ErrFailedConvertID
	}

	return id, nil
}

func (r *UsersRepo) GetByLogin(ctx context.Context, login string) (model.User, error) {
	var user model.User

	filter := bson.M{
		"$or": []bson.M{
			{"username": login},
			{"email": login},
		},
	}

	if err := r.db.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

func (r *UsersRepo) GetSessionInfoByUserID(ctx context.Context, userID bson.ObjectID) (model.UserResponse, error) {
	var user model.UserResponse
	filter := bson.M{"_id": userID}
	if err := r.db.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.UserResponse{}, model.ErrUserNotFound
		}
		return model.UserResponse{}, err
	}
	return user, nil
}
