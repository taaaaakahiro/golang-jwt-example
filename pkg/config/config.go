package config

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Server *serverConfig
	DB     *databaseConfig
	Auth   *authConfig
}

type serverConfig struct {
	Port int `env:"PORT,required"`
}

type databaseConfig struct {
	URI    string `env:"URI,required"`
	Source string `env:"MONGODB_SOURCE,required"`
	Database string `env:"MONGODB_DATABASE,required"`
	
}

type authConfig struct {
	AccessTokenSecret           string `env:"ACCESS_TOKEN_SECRET,default=access_token_secret"`
	RefreshTokenSecret          string `env:"REFRESH_TOKEN_SECRET,default=refresh_token_secret"`
	AccessTokenExpiredDuration  int64  `env:"ACCESS_TOKEN_EXPIRED_DURATION,default=3600000000000"`  // time.Hour * 1
	RefreshTokenExpiredDuration int64  `env:"REFRESH_TOKEN_EXPIRED_DURATION,default=7200000000000"` // time.Hour * 2
}

func LoadConfig(ctx context.Context) (*Config, error) {
	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *Config) Address() string {
	return fmt.Sprintf(":%d", cfg.Server.Port)
}