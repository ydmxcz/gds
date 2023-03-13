package linkedlist_test

import (
	"fmt"
	"testing"

	"github.com/ydmxcz/gds/collections/linkedlist"
	"github.com/ydmxcz/gds/iterator"
)

func TestLinkedList(t *testing.T) {
	l := linkedlist.New[int]()
	for i := 0; i < 10; i++ {
		l.PushBack(i)
	}
	iter, stop := iterator.CastToPull(l.Iter().All)
	for v, ok := iter(); ok; v, ok = iter() {
		fmt.Println(v)
	}
	stop()
}

func TestLinkedListSplitableIter(t *testing.T) {
	l := linkedlist.New[int]()
	for i := 1; i <= 10; i++ {
		l.PushBack(i)
	}
	iter := l.SplitableIter()(3)
	for iter1, ok := iter(); ok; iter1, ok = iter() {
		for v1, o1 := iter1(); o1; v1, o1 = iter1() {
			fmt.Println(v1)
		}
		fmt.Println("====================")
	}
}

func TestBoardCast(t *testing.T) {
	// map[string]strings.Bui
}
