package redis

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestStreamList_GetJobs(t *testing.T) {
	topic := "new test stream1112222232"
	group := "new-test-stream1122222322212-group"
	consumer := "new-test-stream1222232121-group-consumer"
	consume := NewConsume(topic, group, consumer)
	p := StreamList{}
	go func() {
		var i = 0
		for {
			_ = p.Push(context.TODO(), &Message{
				Topic: topic,
				Key:   "key" + fmt.Sprintf("%d", i),
				Body:  []byte(topic),
			})
			time.Sleep(time.Second * 1)
		}
	}()
	time.Sleep(time.Second)

	for {
		msg, err := consume.GetJobs(context.TODO(), 2)
		if err != nil {
			log.Println("err", err)
		}
		for _, v := range msg {
			log.Println("message info ", v.Values)
			_ = consume.Consume(context.TODO(), v.ID)
		}
	}

}
