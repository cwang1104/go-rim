package redis

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Topic string        // 主题
	Key   string        // key 消息唯一key
	Body  []byte        // 消息体
	Delay time.Duration // 延迟时间 秒
}

const (
	ChatPushTopic = "ChatPush"
)

func NewChatPushMsg(body []byte) *Message {
	return &Message{
		Topic: ChatPushTopic,
		Key:   uuid.NewString(),
		Body:  body,
		Delay: time.Second * 5,
	}
}
