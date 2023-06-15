package ws

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

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

func (s *Server) JoinServer(w http.ResponseWriter, r *http.Request, responseHeader http.Header, uid int64, oids []int64) error {
	websockt, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return err
	}

	conn := NewConn(websockt)
	go conn.Reader()
	go conn.Writer()

	for {
		select {
		case <-conn.Done():
			return err
		case msg := <-conn.ReadFromReadChan():
			s.handleMsg(msg)
		}
	}

	return nil
}

func (s *Server) addConn(uid int64, oids []int64) {

}

func (s *Server) handleMsg(msg *Msg) {
	switch msg.MsgType {
	case Ping:

	}
}
