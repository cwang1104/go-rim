package ws

import (
	"context"
	"encoding/json"
	"log"
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
