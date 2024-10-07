package api

import (
	"auth-service/internal/config"
	"auth-service/pkg/database/mongodb"
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
		logger.Fatal("Failed to connect to database",
			zap.Error(err),
			zap.String("context", "Initializing application"),
			zap.String("version", "1.0.0"),
			zap.String("environment", "development"),
		)
	}
	db := mongoClient.Database(cfg.Mongo.Name)
	repositories := repo.NewRepositories(db)
	services := service.NewServices(repositories)
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

	if err := db.Close(); err != nil {
		logger.Errorf("failed to stop postgres database: %v", err)
	}
}
