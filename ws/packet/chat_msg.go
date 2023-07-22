package packet

import "time"

type ChatMsgType uint8

const (
	Text ChatMsgType = iota + 1
	Img
	Video
	Ack
)

type MsgReadStatus uint8

const (
	UnRead  MsgReadStatus = iota
	HasRead MsgReadStatus = iota + 1
)

const (
	TypeYes uint8 = 1
	TypeNo  uint8 = 0
)

// SentChatMsg 单聊、发送的消息
type SentChatMsg struct {
	Text      string      `json:"text"`
	File      *FileInfo   `json:"file"`
	ReceiveID int64       `json:"receiveID"`
	Type      ChatMsgType `json:"type"`

	SenderID  string `json:"senderID"`
	Timestamp int64  `json:"timestamp"`
}

// ChatMsg 回执、通知类、收到的聊天的消息
type ChatMsg struct {
	SessionId  int64         `json:"sessionId"`
	SeqID      int64         `json:"seqID"`      //消息排序id
	MsgID      int64         `json:"msgID"`      //消息id
	FileInfo   *FileInfo     `json:"fileInfo"`   //图片或者视频消息信息
	IsRead     MsgReadStatus `json:"isRead"`     //是否已读
	IsReceive  uint8         `json:"isReceive"`  //是否收到
	IsPeerRead MsgReadStatus `json:"isPeerRead"` //自己已读
	MsgType    ChatMsgType   `json:"msgType"`    //消息类型
	Text       string        `json:"text"`       //消息文本内容
	Group      int64         `json:"group"`      //消息分组
	Sender     *SenderInfo   `json:"sender"`     //发送者信息
	ReceiveID  int64         `json:"receiveID"`  //接收者id
	Timestamp  int64         `json:"createTime"` //创建时间
}

type FileInfo struct {
	FileUrl     string `json:"fileUrl"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
}

type SenderInfo struct {
	UserID           int64
	Avatar, Nickname string
}

func NewMsgAck(toUserID int64, isRead MsgReadStatus, sender *SenderInfo, msgID int64) *ChatMsg {
	return &ChatMsg{
		ReceiveID: toUserID,
		MsgType:   Ack,
		IsRead:    isRead,
		Sender:    sender,
		MsgID:     msgID,
		Timestamp: time.Now().Unix(),
	}
}
