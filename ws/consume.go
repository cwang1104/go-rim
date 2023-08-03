package ws

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"ws/pkg/db/redis"
	"ws/ws/packet"
)

func ConsumeChatPushMsg() {
	t := redis.NewChatPushConsume()
	for {
		msg, _ := t.GetJobs(context.TODO(), 2)
		for _, message := range msg {
			for _, val := range message.Values {
				data, ok := val.(string)
				if !ok {
					continue
				}
				var chatMsg packet.SentChatMsg
				err := json.Unmarshal([]byte(data), &chatMsg)
				if err != nil {
					log.Println("json unmarshal err", err, "str", data)
					continue
				}

				conn, err := WsServer.getConn(chatMsg.ReceiveID)
				if err != nil {
					log.Printf("user %d is offline", chatMsg.ReceiveID)
					continue
				}

				sendMsg := packet.NewV1Msg(packet.Chat)
				sendMsg.Content = chatMsg

				err = conn.SendToWriteChan(context.Background(), sendMsg)
				if err != nil {
					log.Printf("send to failed,err = %v", err)
					continue
				}
			}
		}
	}
}

func ConsumeDelayList() {
	delay := redis.DelayQueue{}

	listMsg := redis.NewChatPushMsg(nil)
	for {
		keys := delay.GetJobKeys(context.Background(), listMsg.Topic)
		for _, v := range keys {
			data := delay.GetJob(context.Background(), listMsg.Topic, v)
			if data == "" {
				continue
			}
			var chatMsg packet.SentChatMsg
			err := json.Unmarshal([]byte(data), &chatMsg)
			if err != nil {
				log.Println("json unmarshal err", err, "str", data)
				continue
			}

			ack := redis.GetSendAckKey(chatMsg.MsgID)
			if ack == 1 {
				sList := redis.StreamList{}
				listMsg.Body = []byte(data)
				_ = sList.Push(context.Background(), listMsg)
			}
			delay.Consume(context.Background(), listMsg.Topic, v)
		}
		time.Sleep(time.Second)
	}
}
