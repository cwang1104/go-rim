package redis

import (
	"context"
	"fmt"
)

type OfflineList struct {
	Uid int64
}

func (o *OfflineList) Push(ctx context.Context, body []byte) error {
	return client.LPush(ctx, o.getKey(), body).Err()
}

func (o *OfflineList) GetOne(ctx context.Context) string {
	return client.RPop(ctx, o.getKey()).Val()
}

func (o *OfflineList) getKey() string {
	return fmt.Sprintf("offline:uid-%d", o.Uid)
}
