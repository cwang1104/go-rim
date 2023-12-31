package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	client   *redis.Client
	rootCtx  = context.TODO()
	baseTime = time.Second * 2
)

func Init() error {
	var err error
	client, err = New()
	return err
}

func GetClient() *redis.Client {
	return client
}

func New() (*redis.Client, error) {
	opt := redis.Options{
		Addr:            "1",
		DB:              1,
		Username:        "",
		Password:        "111",
		PoolSize:        10,
		MinIdleConns:    256,
		ConnMaxIdleTime: time.Second * time.Duration(120),
		ConnMaxLifetime: time.Second * time.Duration(120),
	}

	client = redis.NewClient(&opt)
	ctx, cancel := context.WithTimeout(rootCtx, baseTime*3)
	defer cancel()

	return client, client.Ping(ctx).Err()
}
