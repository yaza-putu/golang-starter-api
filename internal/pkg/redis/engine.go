package redis_client

import "github.com/redis/go-redis/v9"

var Instance *redis.Client

// Mock redis connection and action
func Mock(r *redis.Client) {
	Instance = r
}
