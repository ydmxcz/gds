package stream_test

import (
	"fmt"
	"testing"

	"github.com/ydmxcz/gds/collections/slice"
	"github.com/ydmxcz/gds/stream"
)

func TestTake(t *testing.T) {
	stm := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Stream(4)

	stream.Collect(stream.Take(stm, 5), func(a int) {
		fmt.Println(a)
	})

}
