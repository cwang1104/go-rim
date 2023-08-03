package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"time"
	"ws/ws/packet"
)

func GetMessage() {
	url := "ws://127.0.0.1:60000/api/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("dial" + err.Error())
		return
	}
	defer conn.Close()

	go func() {
		loginMsg := packet.NewV1Msg(packet.Auth)
		loginMsg.Content = packet.AuthMsg{
			Token: "token2",
		}

		err = conn.WriteJSON(loginMsg)
		if err != nil {
			log.Println("send auth failed", err)
			panic(err)
		}
		i := 0
		for {
			pingMsg := packet.NewV1Msg(packet.Ping)
			pingMsg.Content = packet.PingMessage{
				Text: "ping",
			}
			//err := conn.WriteMessage(websocket.BinaryMessage, []byte("ping"+fmt.Sprintf("%d", i)))
			err = conn.WriteJSON(pingMsg)
			//log.Println("send ping")
			if err != nil {
				log.Println(err)
				return
			}
			time.Sleep(time.Second * 1)

			i++
			sendMsg := packet.NewV1Msg(packet.Chat)
			sendMsg.Content = packet.SentChatMsg{
				MsgID:     uuid.NewString(),
				Text:      fmt.Sprintf("message %d", i),
				ReceiveID: 1,
				Type:      packet.Text,
				SenderID:  2,
				Timestamp: time.Now().Unix(),
			}

			err = conn.WriteJSON(sendMsg)
			if err != nil {
				log.Println(err)
				return
			}
			if i == 50 {
				return
			}

		}
	}()

	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("read msg error " + err.Error())
			return
		}
		var msg packet.Msg
		_ = json.Unmarshal(data, &msg)
		if msg.MsgType != packet.Pong {
			log.Println("msg", msg.Content)
		}

	}
}

func main() {
	go GetMessage()
	a := make(chan struct{})
	for {
		fmt.Println("start")
		xx := <-a
		log.Println(xx)

	}
}
