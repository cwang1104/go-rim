package redis

import (
	"context"
	"fmt"
	"log"
	"time"
)

func SetSendAckKey(msgId string) {
	err := client.Set(context.TODO(), getSendExKey(msgId), 1, time.Second*15).Err()
	if err != nil {
		log.Println("SetSendAckKey err", err, "msgId", msgId)
	}
}

func DelSendAckKey(msgId string) {
	key, err := client.GetDel(context.Background(), getSendExKey(msgId)).Result()
	if err != nil {
		log.Println("DelSendAckKey err", err, "msgId", msgId)
	}
	log.Println(getSendExKey(msgId), key)
}

func GetSendAckKey(msgId string) int {
	a, err := client.Get(context.Background(), getSendExKey(msgId)).Int()
	if err != nil {
		log.Println("GetSendAckKey err", err, "msgId", msgId, "a", a)
	}
	return a
}

func getSendExKey(msgId string) string {
	return fmt.Sprintf("send-ack-key-%s", msgId)
}
