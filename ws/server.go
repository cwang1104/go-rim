package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"ws/pkg/util"
	"ws/ws/packet"
	"ws/ws/socket"

	"github.com/gorilla/websocket"
)

const timeout = time.Second * 20

var upgrader = websocket.Upgrader{

	//校验请求来源
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	context context.Context

	container map[int64]*socket.Conn
	mux       sync.RWMutex
}

var WsServer *Server

func init() {
	WsServer = &Server{
		context:   context.TODO(),
		container: make(map[int64]*socket.Conn),
	}
}

func (s *Server) JoinServer(w http.ResponseWriter, r *http.Request, responseHeader http.Header) {
	websockt, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		log.Println("upgrade ws failed:", err)
		return
	}
	s.proc(websockt)
}

func (s *Server) login(wsConn *websocket.Conn) (uid int64, err error) {
	_, data, err := wsConn.ReadMessage()
	if err != nil {
		return 0, err
	}

	authMsg := packet.NewV1Msg(packet.Auth)
	authMsg.Content = packet.AuthMsg{}
	err = json.Unmarshal(data, authMsg)

	if authMsg.MsgType != packet.Auth {
		return 0, errors.New("need auth first")
	}

	token, err := packet.ContentToStruct[packet.AuthMsg](authMsg.Content)
	if err != nil {
		return 0, err
	}
	//todo: 验证token正确，返回uid
	uid, err = util.VerifyToken(token.Token)
	if err != nil {
		return 0, err
	}
	log.Printf("user %d login success", err)
	return
}

func (s *Server) proc(WebSocketConn *websocket.Conn) {

	uid, err := s.login(WebSocketConn)
	if err != nil {
		return
	}
	conn := socket.NewConn(WebSocketConn)
	defer conn.Close()
	s.addConn(conn, uid)

	go conn.Reader()
	go conn.Writer()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			log.Println("time out ", conn.Info.UserID)
			return
		case <-conn.Done():
			log.Println("done ", conn.Info.UserID)
			return
		case msg := <-conn.ReadFromReadChan():
			timer.Reset(timeout)
			s.handleMsg(conn, msg)
		}
	}

}

func (s *Server) Broadcast(oid int64, msg *packet.Msg) {
	//s.mux.RLock()
	//defer s.mux.RUnlock()
	//for _, v := range s.container {
	//	for _, conn := range v {
	//		_ = conn.SendToWriteChan(context.TODO(), msg)
	//	}
	//}
}

func (s *Server) addConn(conn *socket.Conn, uid int64) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.container[uid] = conn
}

func (s *Server) removeConn(uid int64) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if val, ok := s.container[uid]; ok {
		val.Close()
	}
	log.Printf("user %d login out", uid)
	delete(s.container, uid)
}

func (s *Server) getConn(uid int64) (*socket.Conn, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if cn, ok := s.container[uid]; ok {
		return cn, nil
	}
	return nil, errors.New("empty")
}

func (s *Server) handleMsg(conn *socket.Conn, msg *packet.Msg) {
	var resMsg *packet.Msg
	//body := msg.Content
	switch msg.MsgType {
	case packet.Ping:
		resMsg = packet.NewV1Msg(packet.Pong)
		resMsg.Content = packet.NewPongMessage()
	case packet.Quit:
		s.handleQuitMsg(conn)
	case packet.Chat:
		chatMsg, err := packet.ContentToStruct[packet.SentChatMsg](msg.Content)
		if err != nil {
			return
		}
		chatResp, err := s.handleChatMsg(conn, chatMsg)
		if err != nil {
			return
		}

		resMsg = packet.NewV1Msg(packet.ChatAck)
		resMsg.Content = chatResp

		toRecieveMsg := packet.NewV1Msg(packet.Chat)
		chatTo := packet.SentChatMsg{
			Text:      chatMsg.Text,
			ReceiveID: chatMsg.ReceiveID,
			Type:      packet.Text,
			SenderID:  chatMsg.SenderID,
			Timestamp: chatMsg.Timestamp,
		}
		cnn, err := s.getConn(chatMsg.ReceiveID)
		if err != nil {
			log.Println("not exist", chatMsg.ReceiveID)
		}
		toRecieveMsg.Content = chatTo
		err = cnn.SendToWriteChan(context.Background(), toRecieveMsg)
		if err != nil {
			log.Println("send to write chan failed", err)
		}
		//todo: 改写到消息队列统一推送

	case packet.ChatAck:

	case packet.Push:
		return
	}

	if resMsg == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := conn.SendToWriteChan(ctx, resMsg); err != nil {
		log.Println("SendToWriteChan failed: ", conn.Info)
		return
	}
}

func (s *Server) handleQuitMsg(conn *socket.Conn) {
	s.removeConn(conn.Info.UserID)
}

func (s *Server) handleChatMsg(conn *socket.Conn, msg *packet.SentChatMsg) (*packet.ChatMsg, error) {
	msg.Timestamp = time.Now().Unix()

	respMsg := packet.ChatMsg{
		SessionId:  0,
		MsgID:      0,
		SeqID:      0,
		IsRead:     packet.HasRead,
		IsPeerRead: packet.UnRead,
		//IsReceive:  packet.TypeYes,
		ReceiveID: msg.ReceiveID,
		MsgType:   packet.Ack,
		Sender: &packet.SenderInfo{
			UserID: conn.Info.UserID,
		},
		Text:      msg.Text,
		Timestamp: msg.Timestamp,
	}

	return &respMsg, nil
}

func (s *Server) Close(uid int64) {
	s.removeConn(uid)
	fmt.Println(s.container)
}
