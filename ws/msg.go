package ws

type Msg struct {
	Version uint16
	MsgType uint16
	Content any
}
