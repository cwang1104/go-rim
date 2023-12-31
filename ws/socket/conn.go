package socket

import (
	"context"
	"errors"
	"sync/atomic"
	"ws/ws/packet"

	"github.com/gorilla/websocket"
)

var ErrorConnClosed = errors.New("conn is already closed.")

type ConnInfo struct {
	UserID int64
}

type Conn struct {
	websocket           *websocket.Conn
	context             context.Context
	writeChan, readChan chan *packet.Msg

	Info   ConnInfo
	signal chan struct{}
	closed int32
}

func NewConn(conn *websocket.Conn) *Conn {
	return &Conn{
		websocket: conn,
		context:   context.TODO(),
		writeChan: make(chan *packet.Msg, 1<<7),
		readChan:  make(chan *packet.Msg, 1<<7),
		signal:    make(chan struct{}),
		Info:      ConnInfo{},
		closed:    0,
	}
}

func (c *Conn) SetConnInfo(uid int64) *Conn {
	c.Info = ConnInfo{
		UserID: uid,
	}
	return c
}

func (c *Conn) CheckClosed() bool {
	if atomic.LoadInt32(&c.closed) == 1 {
		return true
	} else {
		return false
	}
}

func (c *Conn) SendToWriteChan(ctx context.Context, msg *packet.Msg) error {
	if c.CheckClosed() {
		return ErrorConnClosed
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.writeChan <- msg:
	}
	return nil
}

func (c *Conn) ReadFromReadChan() <-chan *packet.Msg {
	return c.readChan
}

func (c *Conn) Reader() error {
	for {
		msg, err := packet.NormalRead(c.websocket)
		//log.Println("read: ", msg)
		if err != nil {
			c.SetClose()
			close(c.signal)
			return err
		}
		c.readChan <- msg
	}
}

func (c *Conn) Writer() error {
	for {
		select {
		case msg := <-c.writeChan:
			err := packet.Write(c.websocket, msg)
			if err != nil {
				return err
			}
		case <-c.signal:
			return ErrorConnClosed
		}
	}
}

func (c *Conn) SetClose() {
	atomic.StoreInt32(&c.closed, 1)
	c.Close()
}

func (c *Conn) Close() error {
	return c.websocket.Close()
}

func (c *Conn) Done() <-chan struct{} {
	return c.signal
}
