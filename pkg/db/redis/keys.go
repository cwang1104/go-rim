package redis

import (
	"context"
	"fmt"
	"time"
)

func SetSendAckKey(msgId string) {
	client.Set(context.TODO(), getSendExKey(msgId), 1, time.Second*15)
}

func DelSendAckKey(msgId string) {
	client.GetDel(context.Background(), getSendExKey(msgId))
}

func GetSendAckKey(msgId string) int {
	a, _ := client.Get(context.Background(), getSendExKey(msgId)).Int()

	return a
}

func getSendExKey(msgId string) string {
	return fmt.Sprintf("send-ack-key-%s", msgId)
}
