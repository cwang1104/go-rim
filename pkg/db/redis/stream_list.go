package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

const (
	ExpireTime3  = time.Second * 3
	ExpireTime5  = time.Second * 5
	ExpireTime10 = time.Second * 10
)

type StreamList struct {
}

const (
	MaxListLen = 50000
)

func (s *StreamList) Push(ctx context.Context, msg *Message) error {
	err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: msg.Topic,
		ID:     "*",
		Values: map[string]interface{}{
			msg.Key: msg.Body,
		},
		MaxLen: MaxListLen,
	}).Err()

	if err != nil {
		return err
	}
	return nil
}

type StreamConsume struct {
	GroupName string
	Consumer  string
	Topic     string
}

func NewConsume(topic, group, consumer string) *StreamConsume {
	s := &StreamConsume{
		GroupName: group,
		Consumer:  consumer,
		Topic:     topic,
	}
	ctx, cancel := context.WithTimeout(context.TODO(), ExpireTime3)
	defer cancel()

	_ = s.createConsumeGroup(ctx)
	_ = s.createConsumer(ctx)
	return s
}

func (s *StreamConsume) GetJobs(ctx context.Context, count int64) ([]redis.XMessage, error) {
	msg, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    s.GroupName,
		Consumer: s.Consumer,
		Streams:  []string{s.Topic, ">"},
		Count:    count,
	}).Result()
	if err != nil && err != redis.Nil {
		log.Println("XReadGroup err", err)
		return nil, err
	}
	if msg != nil {
		return msg[0].Messages, nil
	}
	return nil, nil
}

func (s *StreamConsume) Consume(ctx context.Context, id string) error {
	return client.XAck(ctx, s.Topic, s.GroupName, id).Err()
}

func (s *StreamConsume) createConsumeGroup(ctx context.Context) error {
	return client.XGroupCreateMkStream(ctx, s.Topic, s.GroupName, "0").Err()
}

func (s *StreamConsume) createConsumer(ctx context.Context) error {
	return client.XGroupCreateConsumer(ctx, s.Topic, s.GroupName, s.Consumer).Err()
}
