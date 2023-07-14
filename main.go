package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
	"ws/ws"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	r := gin.Default()

	r.GET("/api/ws", func(c *gin.Context) {
		ws.WsServer.JoinServer(c.Writer, c.Request, nil)
	})

	r.Run(":60000")

}

func id() int64 {
	return rand.Int63n(5000)
}
