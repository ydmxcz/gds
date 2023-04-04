package delayqueue

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ydmxcz/gds/collections/priorityqueue"
	"github.com/ydmxcz/gds/fn"
)

type Delayed interface {
	Expiration() int64
}

type delayd[T Delayed] struct {
	elem       T
	expiration int64
}

func (d *delayd[T]) Expiration() int64 {
	return d.expiration
}

// 延迟队列
type DelayQueue[T Delayed] struct {
	mutex    sync.Mutex                       // 互斥锁
	pq       *priorityqueue.Queue[*delayd[T]] // 优先队列
	comp     fn.Compare[T]
	sleeping int32         // 已休眠
	wakeupC  chan struct{} // 唤醒队列的通知
}

func NewDelayQueue[T Delayed](comp fn.Compare[T], size int) *DelayQueue[T] {

	return &DelayQueue[T]{
		pq: priorityqueue.New(func(d1, d2 *delayd[T]) int {
			if d1.expiration < d2.expiration {
				return -1
			} else if d1.expiration > d2.expiration {
				return 1
			} else {
				return 0
			}
		}, size), // 优先队列
		comp:    comp,
		wakeupC: make(chan struct{}), // 无缓冲管道saw
	}
}

// 添加元素到队列
func (dq *DelayQueue[T]) Push(elem T) bool {

	dq.mutex.Lock()
	dq.pq.Push(&delayd[T]{
		elem:       elem,
		expiration: elem.Expiration(),
	})
	d, _ := dq.pq.Peek()
	dq.mutex.Unlock()

	if dq.comp(d.elem, elem) == 0 {
		// 如果延迟队列为休眠状态，唤醒他
		if atomic.CompareAndSwapInt32(&dq.sleeping, 1, 0) {
			// 唤醒可能会发生阻塞
			dq.wakeupC <- struct{}{}
		}
	}
	return true
}

// 断地等待一个元素过期，然后将过期的元素发送到通道C。
func (dq *DelayQueue[T]) PopWithCtx(ctx context.Context) (val T, ok bool) {
	for {

		var delta int64 = 0
		dq.mutex.Lock()
		item, ok := dq.pq.Peek()
		if ok {
			now := time.Now().UnixMilli()
			if item.expiration > now {
				delta = item.expiration - now
				item = nil
			} else {
				dq.pq.Pop()
			}
		}

		dq.mutex.Unlock()
		if item == nil {
			//没有要过期的定时器，	将延迟队列设置为休眠
			//为什么要用atomic原子函数，是为了防止Offer 和 Poll出现竞争
			atomic.StoreInt32(&dq.sleeping, 1)
		}

		if item == nil {
			if delta == 0 {
				// 说明延迟队列中已经没有timer，因此等待新的timer添加时wake up通知，或者等待退出通知
				select {
				case <-dq.wakeupC:
					continue
				case <-ctx.Done():
					goto exit
				}
			} else if delta > 0 {
				// 说明延迟队列中存在未过期的定时器
				select {
				case <-dq.wakeupC:
					// 当前定时器已经是休眠状态，如果添加了一个比延迟队列中最早过期的定时器更早的定时器,延迟队列被唤醒
					continue
				case <-time.After(time.Duration(delta) * time.Millisecond):
					// timer.After添加了一个相对时间定时器,并等待到期

					if atomic.SwapInt32(&dq.sleeping, 0) == 0 {
						//防止被阻塞
						<-dq.wakeupC
					}
					continue
				case <-ctx.Done():
					goto exit
				}
			}
		}

		select {
		// case dq.C <- item.node:
		case <-ctx.Done():
			goto exit
		default:
			return item.elem, true
		}
	}

exit:
	// Reset the states
	atomic.StoreInt32(&dq.sleeping, 0)
	return
}

// Poll启动一个无限循环，在这个循环中它不断地等待一个元素过期，然后将过期的元素发送到通道C。
func (dq *DelayQueue[T]) Pop() (T, bool) {
	for {

		var delta int64 = 0
		dq.mutex.Lock()
		item, ok := dq.pq.Peek()
		if ok {
			now := time.Now().UnixMilli()
			if item.expiration > now {
				delta = item.expiration - now
				item = nil
			} else {
				dq.pq.Pop()
			}

		}

		dq.mutex.Unlock()
		if item == nil {
			//没有要过期的定时器，	将延迟队列设置为休眠
			//为什么要用atomic原子函数，是为了防止Offer 和 Poll出现竞争
			atomic.StoreInt32(&dq.sleeping, 1)
		}

		if item == nil {
			if delta == 0 {
				// 说明延迟队列中已经没有timer，因此等待新的timer添加时wake up通知，或者等待退出通知
				select {
				case <-dq.wakeupC:
					continue
				}
			} else if delta > 0 {
				// 说明延迟队列中存在未过期的定时器
				select {
				case <-dq.wakeupC:
					// 当前定时器已经是休眠状态，如果添加了一个比延迟队列中最早过期的定时器更早的定时器,延迟队列被唤醒
				case <-time.After(time.Duration(delta) * time.Millisecond):
					// timer.After添加了一个相对时间定时器,并等待到期
					if atomic.SwapInt32(&dq.sleeping, 0) == 0 {
						//防止被阻塞
						<-dq.wakeupC
					}
				}
				continue
			}

		}
		return item.elem, true
	}
}

func (dq *DelayQueue[T]) Poll() (val T, ok bool) {

	dq.mutex.Lock()
	defer dq.mutex.Unlock()

	item, ok := dq.pq.Peek()
	if ok {
		if (item.expiration - time.Now().UnixMilli()) > 0 {
			return
		} else {
			return item.elem, true
		}
	}
	return
}

func (dq *DelayQueue[T]) Len() int {
	return dq.pq.Len()
}
