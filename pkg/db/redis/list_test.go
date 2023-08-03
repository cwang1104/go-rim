package redis

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestStreamList_GetJobs(t *testing.T) {
	tp := "tp222--stream"
	q := StreamList{}
	group := "group22"
	consumer := "group22-consumer1"
	for i := 0; i < 10; i++ {
		msg := Message{
			Topic: tp,
			Key:   fmt.Sprintf("st1=keys+%d", i),
			Body:  []byte("test1 st queue " + fmt.Sprintf("keys+%d", i)),
			Delay: time.Second * 2,
		}
		err := q.Push(context.TODO(), &msg)
		require.NoError(t, err)
	}

	//err := q.CreateConsumeGroup(context.TODO(), tp, group)
	//require.NoError(t, err)
	//
	//err = q.CreateConsumer(context.TODO(), tp, group, consumer)
	//require.NoError(t, err)

	for i := 0; i < 20; i++ {

		msg, err := q.GetJobs(context.TODO(), tp, group, consumer)
		require.NoError(t, err)
		log.Println("msg", msg)
		for _, v := range msg {
			err = q.Consume(context.TODO(), tp, v.ID, group)
			require.NoError(t, err)
		}
		time.Sleep(time.Second)
	}
}
