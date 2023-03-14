package delayqueue_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ydmxcz/gds/collections/delayqueue"
)

type myObj struct {
	exp  int64
	name string
}

func (mo *myObj) Expiration() int64 {
	return mo.exp
}
func TestDelayQueue(t *testing.T) {
	dq := delayqueue.NewDelayQueue[*myObj](10)
	dq.Push(&myObj{
		name: "mcz1",
		exp:  time.Now().Add(time.Second * 1).UnixMilli(),
	})
	dq.Push(&myObj{
		name: "mcz2",
		exp:  time.Now().Add(time.Second * 2).UnixMilli(),
	})
	dq.Push(&myObj{
		name: "mcz3",
		exp:  time.Now().Add(time.Second * 4).UnixMilli(),
	})
	dq.Push(&myObj{
		name: "mcz4",
		exp:  time.Now().Add(time.Second * 5).UnixMilli(),
	})
	for dq.Len() > 0 {
		fmt.Println(dq.Pop())
	}
}
