package packet

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type MsgType uint8

const (
	Ping MsgType = iota + 1
	Pong
	Quit
	Auth

	Chat
	ChatAck
	Push
)

const MsgVersionV1 uint8 = 1

type ReadMsgType interface {
	PingMessage | AuthMsg | SentChatMsg | ChatMsg
}

type ReadMsg[T ReadMsgType] struct {
	Version uint8   `json:"version"`
	MsgType MsgType `json:"msgType"`
	Content T       `json:"content"`
}

type Msg struct {
	Version uint8   `json:"version"`
	MsgType MsgType `json:"msgType"`
	Content any     `json:"content"`
}

func ContentToStruct[T ReadMsgType](content any) (*T, error) {
	data, _ := json.Marshal(content)
	var info T
	err := json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
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
