package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type MsgVersion uint8

const (
	Ping MsgVersion = iota + 1
	Quit
)

type Msg struct {
	Version uint16
	MsgType MsgVersion
	Content any
}

func Read(conn *websocket.Conn) (*Msg, error) {

	_, data, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	var msg Msg
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func Write(conn *websocket.Conn, msg *Msg) error {
	return conn.WriteJSON(msg)
}
