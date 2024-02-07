package core

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/yaza-putu/golang-starter-api/internal/config"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/redis"
)

func Redis() {
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis().Host, config.Redis().Port),
		Password: config.Redis().Password,
		DB:       config.Redis().DB,
	})

	redis_client.Instance = c
}
