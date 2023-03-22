package delayqueue_test

import (
	"fmt"
	"sync"
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
	dq := delayqueue.NewDelayQueue(func(mo1, mo2 *myObj) int {
		if mo1 == mo2 {
			return 0
		}
		return 1
	}, 10)
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

func TestQueue(t *testing.T) {
	q := delayqueue.NewDelayQueue(func(mo1, mo2 *myObj) int {
		if mo1 == mo2 {
			return 0
		}
		return 1
	}, 10)
	// q.Push(&aa{
	// 	val:        888888,
	// 	expiration: time.Now().Add(time.Second * 6).UnixMilli(),
	// })
	q.Push(&myObj{
		name: "sssss",
		exp:  time.Now().Add(time.Second * 8).UnixMilli(),
	})
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		time.Sleep(time.Second)
		q.Push(&myObj{
			name: "aaaa",
			exp:  time.Now().Add(time.Second).UnixMilli(),
		})
		wg.Done()
	}()
	go func() {
		// a, b := q.Pop()
		fmt.Println(q.Pop())
		wg.Done()
	}()
	wg.Wait()
}
