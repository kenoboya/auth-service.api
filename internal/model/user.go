package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	RoleAdmin    = "admin"
	RoleCustomer = "customer"
	RoleShipper  = "shipper"
)

type UserResponse struct {
	UserID bson.ObjectID `bson:"user_id"`
	Role   string        `bson:"role"`
}

type User struct {
	UserID       bson.ObjectID `json:"user_id" bson:"user_id"`
	Username     string        `json:"username" bson:"username"`
	Email        string        `json:"email" bson:"email"`
	Password     string        `json:"password" bson:"password"`
	RegisteredAt time.Time     `json:"registered_at" bson:"registered_at"`
	Role         string        `json:"role" bson:"role"`
	Person
}

type UserSignUp struct {
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Person   Person `json:"Person"`
}
type UserSignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
