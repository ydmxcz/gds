package stream_test

import (
	"fmt"
	"testing"

	"github.com/ydmxcz/gds/collections/slice"
	"github.com/ydmxcz/gds/collections/truple"
	"github.com/ydmxcz/gds/stream"
)

func TestZip(t *testing.T) {
	stm := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Stream(3)
	stm2 := slice.Of(111, 222, 333, 444, 555, 666, 777, 888, 999).Stream()

	stream.Collect(
		stream.Zip(stm, stm2),
		func(a truple.KV[int, int]) {
			fmt.Println(a.Key, a.Val)
		})

}
