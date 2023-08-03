package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

type Message struct {
	Topic string        // 主题
	Key   string        // key 消息唯一key
	Body  []byte        // 消息体
	Delay time.Duration // 延迟时间 秒
}

type DelayQueue struct{}

func (d *DelayQueue) Push(ctx context.Context, msg *Message) error {
	readyTime := time.Now().Add(msg.Delay).UnixMilli()
	pip := client.TxPipeline()
	pip.ZAdd(ctx, d.topicZSet(msg.Topic), redis.Z{
		Score:  float64(readyTime),
		Member: msg.Key,
	})

	pip.HSet(ctx, d.topicHash(msg.Topic), msg.Key, msg.Body)

	cmds, err := pip.Exec(ctx)
	if err != nil {
		return err
	}
	for _, v := range cmds {
		if v.Err() != nil {
			log.Println("failed ", v.Err(), " name ", v.Name())
		}
	}
	return nil
}

func (d *DelayQueue) GetJobKeys(ctx context.Context, topic string) []string {
	mm := time.Now().UnixMilli()
	keys := client.ZRangeByScore(ctx, d.topicZSet(topic), &redis.ZRangeBy{
		Min: strconv.FormatInt(mm-999, 10),
		Max: strconv.FormatInt(mm, 10),
	}).Val()

	return keys
}

func (d *DelayQueue) Consume(ctx context.Context, topic, key string) {

	pip := client.TxPipeline()

	pip.ZRem(ctx, d.topicZSet(topic), key)

	pip.HDel(ctx, d.topicHash(topic), key)

	_, err := pip.Exec(ctx)
	if err != nil {
		log.Println("Consume err", err)
	}
}

func (d *DelayQueue) topicZSet(topic string) string {
	return topic + ":ZSet"
}

func (d *DelayQueue) topicHash(topic string) string {
	return topic + ":Hash"
}
