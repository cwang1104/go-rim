package redis

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func init() {
	New()
}

func TestDelayQueue_GetJobs(t *testing.T) {

	q := DelayQueue{}
	topic := "testTopic"

	for i := 0; i < 10; i++ {
		msg := Message{
			Topic: topic,
			Key:   fmt.Sprintf("keys+%d", i),
			Body:  []byte("test queue " + fmt.Sprintf("keys+%d", i)),
			Delay: time.Second * 2,
		}
		err := q.Push(context.TODO(), &msg)
		require.NoError(t, err)
	}

	for i := 0; i < 12; i++ {
		keys := q.GetJobKeys(context.TODO(), topic)
		log.Println("keys", keys)
		//for _, v := range keys {
		//	q.Consume(context.TODO(), topic, v)
		//}
		time.Sleep(time.Second)
	}
}
