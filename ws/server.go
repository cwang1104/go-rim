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

func (s *Server) JoinServer(w http.ResponseWriter, r *http.Request, responseHeader http.Header) error {
	websockt, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return err
	}
	return s.login(websockt)

}

func (s *Server) login(WebSocketConn *websocket.Conn) error {
	_, data, err := WebSocketConn.ReadMessage()
	if err != nil {
		return err
	}
	authMsg := packet.NewV1Msg(packet.Auth)
	authMsg.Content = packet.AuthMsg{}
	err = json.Unmarshal(data, authMsg)
	if err != nil {
		return err
	}
	//todo：验证token正确
	if true {
		var uid int64 = 10

		defer s.Close(uid)
		conn := socket.NewConn(WebSocketConn)
		conn.SetConnInfo(uid)

		s.addConn(conn, uid)
		log.Println("uid: ", uid)

		go conn.Reader()
		go conn.Writer()

		timer := time.NewTimer(timeout)
		defer timer.Stop()

		for {
			select {
			case <-timer.C:
				return errors.New("time out")
			case <-conn.Done():
				return err
			case msg := <-conn.ReadFromReadChan():
				timer.Reset(timeout)
				s.handleMsg(conn, msg)
			}
		}

		return nil

	} else {
		return nil
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
	//for _, v := range oids {
	//	val, ok := s.container[v]
	//	if !ok {
	//		continue
	//	} else {
	//		delete(val, uid)
	//	}
	//}
}

func (s *Server) handleMsg(conn *socket.Conn, msg *packet.Msg) {
	var resMsg *packet.Msg
	//body := msg.Content
	switch msg.MsgType {
	case packet.Ping:
		resMsg = packet.NewV1Msg(packet.Pong)
		resMsg.Content = packet.NewPongMessage()
	case packet.Quit:

		s.removeConn(conn.Info.UserID)
	case packet.Chat:
		chatMsg, err := packet.ContentToStruct[packet.SentChatMsg](msg.Content)
		if err != nil {
			return
		}
		fmt.Println(chatMsg)
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

func (s *Server) Close(uid int64) {
	s.removeConn(uid)
	fmt.Println(s.container)
}
