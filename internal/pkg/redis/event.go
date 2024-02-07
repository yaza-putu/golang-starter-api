package redis_client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/logger"
)

func Set(ctx context.Context, key string, data any) {
	p, err := json.Marshal(data)
	logger.New(err, logger.SetType(logger.ERROR))
	Instance.Set(ctx, key, p, 0)
}

func Del(ctx context.Context, key string) {
	Instance.Del(ctx, key)
}

func FindSet(ctx context.Context, key string, data any) {
	e := Instance.Get(ctx, key)
	if e.Err() != nil {
		Set(ctx, key, data)
	}
}

func Get(ctx context.Context, key string, dest any) error {
	o := Instance.Get(ctx, key)
	if o.Err() != nil {
		logger.New(o.Err(), logger.SetType(logger.ERROR))
		return o.Err()
	}

	r, err := o.Result()
	if err != nil {
		logger.New(err, logger.SetType(logger.ERROR))
		return err
	}

	fmt.Println("get data from redis")

	return json.Unmarshal([]byte(r), dest)
}
