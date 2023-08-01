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
		Addr:            "120.24.45.199:56379",
		DB:              0,
		Username:        "",
		Password:        "123456",
		PoolSize:        0,
		MinIdleConns:    256,
		ConnMaxIdleTime: time.Second * time.Duration(120),
		ConnMaxLifetime: time.Second * time.Duration(0),
	}

	client = redis.NewClient(&opt)
	ctx, cancel := context.WithTimeout(rootCtx, baseTime*3)
	defer cancel()

	return client, client.Ping(ctx).Err()
}
