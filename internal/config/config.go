package config

import (
	"auth-service/internal/model"
	"auth-service/pkg/database/mongodb"
	"auth-service/pkg/database/redis"
	logger "auth-service/pkg/logger/zap"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"go.uber.org/zap"
)

type Config struct {
	HTTP  HttpConfig
	GRPC  GrpcConfig
	Auth  AuthConfig
	Mongo mongodb.MongoConfig
	Redis redis.RedisConfig
}

type HttpConfig struct {
	Addr           string        `mapstructure:"port"`
	ReadTimeout    time.Duration `mapstructure:"readTimeout"`
	WriteTimeout   time.Duration `mapstructure:"writeTimeout"`
	MaxHeaderBytes int           `mapstructure:"maxHeaderBytes"`
}

type GrpcConfig struct {
}

type AuthConfig struct {
	Salt      string
	SecretKey string
	TokenTTL  time.Duration `mapstructure:"tokenTTL"`
}

func Init(configDIR string, envDIR string) (*Config, error) {
	if err := loadViperConfig(configDIR); err != nil {
		return &Config{}, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return &Config{}, err
	}

	if err := loadFromEnv(&cfg, envDIR); err != nil {
		return &Config{}, err
	}

	return &cfg, nil
}

func unmarshal(config *Config) error {
	if err := viper.UnmarshalKey("http", &config.HTTP); err != nil {
		logger.Error("Failed to unmarshal config file",
			zap.String("prefix", "http"),
			zap.Error(err),
		)
		return err
	}

	if err := viper.UnmarshalKey("grpc", &config.GRPC); err != nil {
		logger.Error("Failed to unmarshal config file",
			zap.String("prefix", "grpc"),
			zap.Error(err),
		)
		return err
	}

	if err := viper.UnmarshalKey("mongo", &config.Mongo); err != nil {
		logger.Error("Failed to unmarshal config file",
			zap.String("prefix", "mongo"),
			zap.Error(err),
		)
		return err
	}

	if err := viper.UnmarshalKey("cache", &config.Auth); err != nil {
		logger.Error("Failed to unmarshal config file",
			zap.String("prefix", "cache"),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func loadFromEnv(cfg *Config, envDIR string) error {
	if err := gotenv.Load(envDIR); err != nil {
		logger.Error(
			zap.String("file", ".env"),
			zap.Error(model.ErrNotFoundEnvFile),
		)
		return model.ErrNotFoundEnvFile
	}

	if err := envconfig.Process("MONGO", &cfg.Mongo); err != nil {
		logger.Error("Failed to unmarshal environment file",
			zap.String("prefix", "MONGO"),
			zap.String("file", ".env"),
			zap.Error(err),
		)
		return err
	}

	if err := envconfig.Process("REDIS", &cfg.Mongo); err != nil {
		logger.Error("Failed to unmarshal environment file",
			zap.String("prefix", "REDIS"),
			zap.String("file", ".env"),
			zap.Error(err),
		)
		return err
	}

	cfg.Auth.Salt = os.Getenv("SALT")
	cfg.Auth.SecretKey = os.Getenv("SECRET_KEY")
	return nil
}

func loadViperConfig(path string) error {
	viper.SetConfigName("server")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Error(
				zap.String("file", "server.yaml"),
				zap.String("path", path),
				zap.Error(model.ErrNotFoundConfigFile),
			)
			return model.ErrNotFoundConfigFile
		} else {
			return err
		}
	}
	return viper.MergeInConfig()
}
