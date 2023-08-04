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
				var chatMsg packet.SentChatMsgPush
				err := json.Unmarshal([]byte(data), &chatMsg)
				if err != nil {
					log.Println("json unmarshal err", err, "str", data)
					continue
				}

				offline := redis.OfflineList{
					Uid: chatMsg.ReceiveID,
				}

				sendMsg := packet.NewV1Msg(packet.Chat)
				sendMsg.Content = chatMsg

				conn, err := WsServer.getConn(chatMsg.ReceiveID)
				if err != nil {
					if chatMsg.ReSend {
						_ = offline.Push(context.Background(), []byte(data))
					}
					log.Printf("user %d is offline", chatMsg.ReceiveID)
					continue
				}

				err = conn.SendToWriteChan(context.Background(), sendMsg)
				if err != nil {
					if chatMsg.ReSend {
						_ = offline.Push(context.Background(), []byte(data))
					}
					log.Printf("send to failed,err = %v", err)
					continue
				}

			}
		}
	}
}

func ConsumeDelayList() {
	delay := redis.DelayQueue{}

	delayMsg := redis.NewChatPushDelay(nil)

	for {
		keys := delay.GetJobKeys(context.Background(), delayMsg.Topic)
		for _, v := range keys {
			data := delay.GetJob(context.Background(), delayMsg.Topic, v)
			if data == "" {
				continue
			}
			var chatMsg packet.SentChatMsgPush
			err := json.Unmarshal([]byte(data), &chatMsg)
			if err != nil {
				log.Println("json unmarshal err", err, "str", data)
				continue
			}

			ack := redis.GetSendAckKey(chatMsg.MsgID)
			if ack == 1 {
				streamMsg := redis.NewChatPushMsg(nil)
				chatMsg.ReSend = true
				chatData, _ := json.Marshal(&chatMsg)
				sList := redis.StreamList{}
				streamMsg.Body = chatData
				_ = sList.Push(context.Background(), streamMsg)
			}
			delay.Consume(context.Background(), delayMsg.Topic, v)
		}
		time.Sleep(time.Second)
	}
}
