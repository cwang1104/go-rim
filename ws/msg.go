package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type MsgType uint8

const (
	Ping MsgType = iota + 1
	Pong
	Quit
)

const MsgVersionV1 uint16 = 1

type Msg struct {
	Version uint16  `json:"version"`
	MsgType MsgType `json:"msgType"`
	Content any     `json:"content"`
}

func Read(conn *websocket.Conn) (*Msg, error) {

	_, data, err := conn.ReadMessage()
	log.Println("data: ", string(data))
	if err != nil {
		return nil, err
	}
	var msg Msg
	err = json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	return &msg, nil
}

func Write(conn *websocket.Conn, msg *Msg) error {
	return conn.WriteJSON(msg)
}

func NewV1Msg(msgType MsgType) *Msg {
	return &Msg{
		Version: MsgVersionV1,
		MsgType: msgType,
	}
}
