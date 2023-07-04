package packet

type PingMessage struct {
	Text string `json:"text"`
}

type PongMessage struct {
	Text string `json:"text"`
}

func NewPongMessage() *PongMessage {
	return &PongMessage{
		Text: "pong",
	}
}
