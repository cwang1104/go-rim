package redis

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOfflineList_GetOne(t *testing.T) {
	l := OfflineList{
		Uid: 1345,
	}

	ctx := context.TODO()
	err := l.Push(ctx, []byte("jack Ma"))
	err = l.Push(ctx, []byte("jack a"))
	err = l.Push(ctx, []byte("jack b"))
	require.NoError(t, err)

	for {
		k := l.GetOne(ctx)
		if k == "" {
			break
		}
		fmt.Println("k", k)
	}

}
