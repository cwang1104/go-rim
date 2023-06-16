package ws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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

	container map[int64]map[int64]*Conn
	mux       sync.RWMutex
}

var WsServer *Server

func init() {
	WsServer = &Server{
		context:   context.TODO(),
		container: make(map[int64]map[int64]*Conn),
	}
}

func (s *Server) JoinServer(w http.ResponseWriter, r *http.Request, responseHeader http.Header, uid int64, oids []int64) error {
	websockt, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return err
	}
	defer s.Close(uid, oids)
	conn := NewConn(websockt)
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

func (s *Server) Broadcast(oid int64, msg *Msg) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	for _, v := range s.container {
		for _, conn := range v {
			_ = conn.SendToWriteChan(context.TODO(), msg)
		}
	}
}

func (s *Server) addConn(conn *Conn, uid int64, oids []int64) {
	s.mux.Lock()
	defer s.mux.Unlock()
	for _, v := range oids {
		val, ok := s.container[v]
		if ok {
			val[uid] = conn
		} else {
			m := make(map[int64]*Conn)
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

func (s *Server) handleMsg(conn *Conn, msg *Msg) {
	var resMsg *Msg
	//body := msg.Content
	switch msg.MsgType {
	case Ping:
		resMsg = NewV1Msg(Pong)
		resMsg.Content = NewPongMessage()
	case Quit:
		//todo: 获取oids
		s.removeConn(conn.Info.UserID, 1)
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

func (s *Server) handleQuitMsg(conn *Conn) {
	mysqlOids := make([]int64, 1)
	s.removeConn(conn.Info.UserID, mysqlOids...)
}

func (s *Server) Close(uid int64, oids []int64) {
	s.removeConn(uid, oids...)
	fmt.Println(s.container)
}
