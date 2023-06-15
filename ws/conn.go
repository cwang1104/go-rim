package ws

import (
	"context"
	"errors"
	"net"
	"sync/atomic"
)

var ErrorConnClosed = errors.New("conn is already closed.")

type ConnInfo struct {
	UserID     int64
	CurrentOid int64
}

type Conn struct {
	socket              net.Conn
	context             context.Context
	writeChan, readChan chan *Msg

	Info   ConnInfo
	signal chan struct{}
	closed int32
}

func NewConn(conn net.Conn) *Conn {
	return &Conn{
		socket:    conn,
		context:   context.TODO(),
		writeChan: make(chan *Msg, 1<<7),
		readChan:  make(chan *Msg, 1<<7),
		signal:    make(chan struct{}),
		Info:      ConnInfo{},
		closed:    0,
	}
}

func (c *Conn) SetConnInfo(uid int64, oid int64) *Conn {
	c.Info = ConnInfo{
		UserID:     uid,
		CurrentOid: oid,
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

func (c *Conn) SendToWriteChan(ctx context.Context, msg *Msg) error {
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

func (c *Conn) ReadFromReadChan() *Msg {
	return <-c.readChan
}

func (c *Conn) Read(p []byte) (n int, err error) {
	return c.socket.Read(p)
}

func (c *Conn) Write(b []byte) (n int, err error) {
	return c.socket.Write(b)
}

func (c *Conn) Close() error {
	if c.CheckClosed() {
		return ErrorConnClosed
	}
	c.closed = 1
	c.signal <- struct{}{}
	return c.socket.Close()
}
