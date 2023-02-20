package io

import (
	"github.com/redis/go-redis/v9"
	"golang-jwt-example/pkg/config"
)

type Redis struct {
	Conn *redis.Client
}

func NewRedisClient(cfg *config.Config) *Redis {
	conn := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	return &Redis{
		Conn: conn,
	}
}
