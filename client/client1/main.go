package main

import (
	"encoding/json"
	"fmt"
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
			Token: "token1",
		}

		err = conn.WriteJSON(loginMsg)
		if err != nil {
			log.Println("send auth failed", err)
			panic(err)
		}

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
			time.Sleep(time.Second * 2)
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
		if msg.MsgType == packet.Chat {
			chatMsg, err := packet.ContentToStruct[packet.SentChatMsg](msg.Content)
			if err != nil {
				log.Println("Err", err)
				continue
			}
			log.Println("message", chatMsg.Text)
			log.Println("message-id", chatMsg.MsgID)
			chatAckMsg := packet.NewMsgAck(0, 0, nil, chatMsg.MsgID)
			msg.Content = chatAckMsg
			msg.MsgType = packet.ChatAck
			_ = conn.WriteJSON(msg)
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
