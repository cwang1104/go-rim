package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"ws/pkg/db/redis"
	"ws/ws"
)

func init() {
	err := redis.Init()
	if err != nil {
		log.Println("redis init failed", err)
		panic(err)
	}
}

func main() {
	ws.Init()
	r := gin.Default()

	r.GET("/api/ws", func(c *gin.Context) {
		ws.WsServer.JoinServer(c.Writer, c.Request, nil)
	})

	r.Run(":60000")

}
