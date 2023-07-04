package ws

import (
	"context"
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

	container map[int64]map[int64]*socket.Conn
	mux       sync.RWMutex
}

var WsServer *Server

func init() {
	WsServer = &Server{
		context:   context.TODO(),
		container: make(map[int64]map[int64]*socket.Conn),
	}
}

func (s *Server) JoinServer(w http.ResponseWriter, r *http.Request, responseHeader http.Header, uid int64, oids []int64) error {
	websockt, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return err
	}
	defer s.Close(uid, oids)
	conn := socket.NewConn(websockt)
	conn.SetConnInfo(uid, 0)

	s.addConn(conn, uid, oids)
	log.Println("uid: ", uid)
	log.Println("oids: ", oids)
	go conn.Reader()
	go conn.Writer()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	fmt.Println(s.container)
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
}

func (s *Server) Broadcast(oid int64, msg *packet.Msg) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	for _, v := range s.container {
		for _, conn := range v {
			_ = conn.SendToWriteChan(context.TODO(), msg)
		}
	}
}

func (s *Server) addConn(conn *socket.Conn, uid int64, oids []int64) {
	s.mux.Lock()
	defer s.mux.Unlock()
	for _, v := range oids {
		val, ok := s.container[v]
		if ok {
			val[uid] = conn
		} else {
			m := make(map[int64]*socket.Conn)
			m[uid] = conn
			s.container[v] = m
		}
	}
}

func (s *Server) removeConn(uid int64, oids ...int64) {
	s.mux.Lock()
	defer s.mux.Unlock()
	for _, v := range oids {
		val, ok := s.container[v]
		if !ok {
			continue
		} else {
			delete(val, uid)
		}
	}
}

func (s *Server) handleMsg(conn *socket.Conn, msg *packet.Msg) {
	var resMsg *packet.Msg
	//body := msg.Content
	switch msg.MsgType {
	case packet.Ping:
		resMsg = packet.NewV1Msg(packet.Pong)
		resMsg.Content = packet.NewPongMessage()
	case packet.Quit:
		//todo: 获取oids
		s.removeConn(conn.Info.UserID, 1)
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
	mysqlOids := make([]int64, 1)
	s.removeConn(conn.Info.UserID, mysqlOids...)
}

func (s *Server) Close(uid int64, oids []int64) {
	s.removeConn(uid, oids...)
	fmt.Println(s.container)
}
