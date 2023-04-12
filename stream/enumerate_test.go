package stream_test

import (
	"fmt"
	"testing"

	"github.com/ydmxcz/gds/collections/slice"
	"github.com/ydmxcz/gds/collections/truple"
	"github.com/ydmxcz/gds/stream"
)

func TestEnumerate(t *testing.T) {
	stm := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Stream(4)

	stream.Collect(stream.Enumerate(stm), func(a truple.KV[int, int]) {
		fmt.Println(a.Key, a.Val)
	})

}
