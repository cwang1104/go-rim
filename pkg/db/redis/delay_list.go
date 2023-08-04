package redis

import (
	"context"
	_ "embed"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed scripts/delay_queue_push.lua
var pushLuaScript string

//go:embed scripts/delay_queue_consume.lua
var consumeJobsLuaScript string

type DelayQueue struct{}

func (d *DelayQueue) Push(ctx context.Context, msg *Message) error {
	readyTime := time.Now().Add(msg.Delay).Unix()
	err := client.Eval(ctx, pushLuaScript, []string{d.topicZSet(msg.Topic), d.topicHash(msg.Topic)}, msg.Key, msg.Body, readyTime).Err()
	if err != nil {
		return err
	}

	//pip := client.TxPipeline()
	//pip.ZAdd(ctx, d.topicZSet(msg.Topic), redis.Z{
	//	Score:  float64(readyTime),
	//	Member: msg.Key,
	//})
	//
	//pip.HSet(ctx, d.topicHash(msg.Topic), msg.Key, msg.Body)
	//
	//cmds, err := pip.Exec(ctx)
	//if err != nil {
	//	return err
	//}
	//for _, v := range cmds {
	//	if v.Err() != nil {
	//		log.Println("failed ", v.Err(), " name ", v.Name())
	//	}
	//}
	return nil
}

func (d *DelayQueue) GetJobKeys(ctx context.Context, topic string, count int64) []string {
	keys := client.ZRangeByScore(ctx, d.topicZSet(topic), &redis.ZRangeBy{
		Min:   "-inf",
		Max:   strconv.FormatInt(time.Now().Unix(), 10),
		Count: count,
	}).Val()
	return keys
}

func (d *DelayQueue) GetJob(ctx context.Context, topic, key string) string {
	return client.HGet(ctx, d.topicHash(topic), key).Val()
}

func (d *DelayQueue) Consume(ctx context.Context, topic, key string) {

	err := client.Eval(ctx, consumeJobsLuaScript, []string{d.topicZSet(topic), d.topicHash(topic)}, key).Err()
	if err != nil {
		log.Println("delay consume", err, "lua", consumeJobsLuaScript)
		return
	}
	//pip := client.TxPipeline()
	//
	//pip.ZRem(ctx, d.topicZSet(topic), key)
	//
	//pip.HDel(ctx, d.topicHash(topic), key)
	//
	//_, err := pip.Exec(ctx)
	//if err != nil {
	//	log.Println("Consume err", err)
	//}
}

func (d *DelayQueue) topicZSet(topic string) string {
	return topic + ":ZSet"
}

func (d *DelayQueue) topicHash(topic string) string {
	return topic + ":Hash"
}
