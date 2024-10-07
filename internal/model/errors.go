package model

import "errors"

var (
	ErrNotFoundConfigFile = errors.New("failed to find config file")
	ErrNotFoundEnvFile    = errors.New("failed to load environment file")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrFailedConvertID    = errors.New("failed to convert inserted ID to ObjectID")
	ErrUserNotFound       = errors.New("user doesn't exists")
	ErrInvalidPassword = errors.New("invalid password")
)
