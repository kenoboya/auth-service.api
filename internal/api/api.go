package api

import (
	"auth-service/internal/config"
	repo "auth-service/internal/repository/mongo"
	http_server "auth-service/internal/server/http"
	"auth-service/internal/service"
	"auth-service/internal/transport/rest"
	"auth-service/pkg/database/mongodb"
	"auth-service/pkg/database/redis"
	logger "auth-service/pkg/logger/zap"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func Run(configDIR string, envDIR string) {
	logger.InitLogger()
	cfg, err := config.Init(configDIR, envDIR)
	if err != nil {
		logger.Fatal("Failed to initialize config",
			zap.Error(err),
			zap.String("context", "Initializing application"),
			zap.String("version", "1.0.0"),
			zap.String("environment", "development"),
		)
	}
	mongoClient, err := mongodb.NewClient(cfg.Mongo)
	if err != nil {
		logger.Fatal("Failed to connect to mongo",
			zap.Error(err),
			zap.String("context", "Initializing application"),
			zap.String("version", "1.0.0"),
			zap.String("environment", "development"),
		)
	}
	db := mongoClient.Database(cfg.Mongo.Name)
	redisClient := redis.NewClient(cfg.Redis)
	repositories := repo.NewRepositories(db)
	deps, err := service.NewDeps(repositories, cfg, redisClient)
	if err != nil {
		logger.Fatal("Failed to connect to redis",
			zap.Error(err),
			zap.String("context", "Initializing application"),
			zap.String("version", "1.0.0"),
			zap.String("environment", "development"),
		)
	}
	services := service.NewServices(deps)
	handler := rest.NewHandler(services)
	httpServer := http_server.NewServer(cfg.HTTP, handler)

	go func() {
		if err := httpServer.Run(); err != nil {
			logger.Fatalf("The http server didn't start: %s\n", err)
		}
	}()

	logger.Info("Http server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Errorf("failed to stop http server: %v", err)
	}

	if err := mongoClient.Disconnect(ctx); err != nil {
		logger.Errorf("failed to stop mongo database: %v", err)
	}
	if err := redisClient.Close(); err != nil {
		logger.Errorf("failed to stop redis: %v", err)
	}
}
