package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
	"ws/ws/packet"
)

var sendChan = make(chan []byte, 10000)
var messageMap = make(map[string]struct{})

func GetMessage() {
	url := "ws://127.0.0.1:60000/api/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("dial" + err.Error())
		return
	}
	defer conn.Close()

	loginMsg := packet.NewV1Msg(packet.Auth)
	loginMsg.Content = packet.AuthMsg{
		Token: "token1",
	}

	conn.WriteJSON(loginMsg)

	go func() {

		for {
			pingMsg := packet.NewV1Msg(packet.Ping)
			pingMsg.Content = packet.PingMessage{
				Text: "ping",
			}
			data, _ := json.Marshal(&pingMsg)
			sendChan <- data
			time.Sleep(time.Second * 2)
		}
	}()

	go func(conn *websocket.Conn) {
		for {
			select {
			case data := <-sendChan:
				conn.WriteMessage(websocket.BinaryMessage, data)
			}
		}
	}(conn)

	i := 0
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
			i++

			chatAckMsg := packet.NewMsgAck(0, 0, nil, chatMsg.MsgID)
			msg.Content = chatAckMsg
			msg.MsgType = packet.ChatAck
			messageMap[chatMsg.MsgID] = struct{}{}
			data, _ := json.Marshal(&msg)
			log.Println("message ", chatMsg.Text, " count ", len(messageMap))
			sendChan <- data
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
