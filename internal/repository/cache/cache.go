package cache

import (
	"auth-service/internal/model"
	logger "auth-service/pkg/logger/zap"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type SessionCache interface {
	SetSession(ctx context.Context, session model.Session, userID string) error
	GetSession(ctx context.Context, sessionToken string) (userID string, err error)
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) SetSession(ctx context.Context, session model.Session, userID string) error {
	sessionKey := fmt.Sprintf("session:%s", session.Token)

	_, err := c.client.HSet(ctx, sessionKey, map[string]interface{}{
		"user_id":    userID,
		"expires_at": session.ExpiresAt.Format(time.RFC3339),
	}).Result()

	if err != nil {
		logger.Error(
			model.ErrSetSessionValue.Error(),
			zap.String("session_key", sessionKey),
			zap.String("user_id", userID),
			zap.Time("expires_at", session.ExpiresAt),
			zap.String("cache", "session"),
			zap.Error(err),
		)
		return model.ErrSetSessionValue
	}

	expiration := time.Until(session.ExpiresAt)
	err = c.client.Expire(ctx, sessionKey, expiration).Err()
	if err != nil {
		logger.Error(
			model.ErrSetSessionExpiry.Error(),
			zap.String("session_key", sessionKey),
			zap.String("user_id", userID),
			zap.Time("expires_at", session.ExpiresAt),
			zap.String("cache", "session"),
			zap.Error(err),
		)
		return model.ErrSetSessionExpiry
	}

	return nil
}

func (c *RedisCache) GetSession(ctx context.Context, sessionToken string) (userID string, err error) {
	sessionKey := fmt.Sprintf("session:%s", sessionToken)

	userID, err = c.client.HGet(ctx, sessionKey, "user_id").Result()
	if err == redis.Nil {
		logger.Error(
			model.ErrNotFoundSession.Error(),
			zap.String("session_key", sessionKey),
			zap.String("cache", "session"),
			zap.Error(err),
		)
		return "", model.ErrNotFoundSession
	} else if err != nil {
		logger.Error(
			zap.String("session_key", sessionKey),
			zap.String("cache", "session"),
			zap.Error(err),
		)
		return "", fmt.Errorf("error retrieving user_id from Redis: %w", err)
	}
	return userID, nil
}
