package config

import (
	"fmt"
	"os"
	"time"

	"github.com/go-yaml/yaml"
)

type AppEnv string

const (
	TestEnv AppEnv = "test"
	DevEnv  AppEnv = "dev"
	ProdEnv AppEnv = "prod"
)

const (
	DBEnvKey                 = "DB_URL"
	AccessTokenSecretEnvKey  = "ACCESS_TOKEN_SECRET"
	RefreshTokenSecretEnvKey = "REFRESH_TOKEN_SECRET"
)

const (
	DefaultGRPCPort = "50051"
	DefaultRestPort = "8080"

	DefaultAccessTokenExpirationDuration  = 30 * time.Minute
	DefaultRefreshTokenExpirationDuration = 24 * time.Hour
)

type Config struct {
	AppEnv                         AppEnv        `yaml:"app_env"`
	DBUrl                          string        `yaml:"db_url"`
	RestPort                       string        `yaml:"rest_port"`
	GRPCPort                       string        `yaml:"grpc_port"`
	AccessTokenExpirationDuration  time.Duration `yaml:"access_token_expiration_duration"`
	RefreshTokenExpirationDuration time.Duration `yaml:"refresh_token_expiration_duration"`
}

func InitConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("unable to read path '%s': %w", path, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("unable to unmarshal yaml config: %w", err)
	}

	if config.AppEnv == "" {
		config.AppEnv = DevEnv
	}
	if config.DBUrl == "" {
		config.DBUrl = os.Getenv(DBEnvKey)
	}
	if config.GRPCPort == "" {
		config.GRPCPort = DefaultGRPCPort
	}
	if config.RestPort == "" {
		config.RestPort = DefaultRestPort
	}
	if config.AccessTokenExpirationDuration == 0 {
		config.AccessTokenExpirationDuration = DefaultAccessTokenExpirationDuration
	}
	if config.RefreshTokenExpirationDuration == 0 {
		config.RefreshTokenExpirationDuration = DefaultRefreshTokenExpirationDuration
	}

	return config, nil
}
