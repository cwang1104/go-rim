package packet

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

func Read[T ReadMsgType](conn *websocket.Conn) (*ReadMsg[T], error) {
	var msg ReadMsg[T]
	data, err := read(conn)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(data, &msg)
	return &msg, nil
}

func NormalRead(conn *websocket.Conn) (*Msg, error) {
	data, err := read(conn)
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

func read(conn *websocket.Conn) ([]byte, error) {
	_, data, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return data, nil
}
