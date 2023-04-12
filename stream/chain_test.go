package stream_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ydmxcz/gds/collections/slice"
	"github.com/ydmxcz/gds/stream"
)

func TestChain(t *testing.T) {
	sli := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	sli2 := slice.Of(111, 222, 333, 444, 555, 666, 777, 888, 999, 101010)
	s := sli.Stream()
	s2 := sli2.Stream(3)

	stream.Collect(stream.Chain(s, s2), func(a int) {
		fmt.Println(a)
	})
}

func TestXxx(t *testing.T) {
	_, cancel := context.WithCancel(context.Background())
	cancel()
	cancel()

}
