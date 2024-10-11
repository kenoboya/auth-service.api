package grpc_handler

import (
	"auth-service/internal/model"
	"auth-service/internal/server/grpc/proto"
	"auth-service/internal/service"
	logger "auth-service/pkg/logger/zap"
	"context"

	"go.uber.org/zap"
)

type AuthHandler struct {
	proto.UnimplementedAuthServiceServer
	services service.Services
}

func NewAuthHandler(services service.Services) *AuthHandler {
	return &AuthHandler{services: services}
}

func (h *AuthHandler) Verify(ctx context.Context, req *proto.TokenRequest) (*proto.UserResponse, error) {
	if req.SessionToken == "" {
		logger.Error(
			zap.String("action", "Verify"),
			zap.Error(model.ErrEmptyParam),
		)
		return nil, model.ErrEmptyParam
	}
	userResponse, err := h.services.Users.Verify(ctx, req.SessionToken)
	if err != nil {
		logger.Error(
			zap.String("action", "Verify"),
			zap.Error(err),
		)
	}

	return &proto.UserResponse{
		UserId: userResponse.UserID.String(),
		Role:   userResponse.Role,
	}, nil
}
