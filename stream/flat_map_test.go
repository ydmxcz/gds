package stream_test

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/ydmxcz/gds/collections/slice"
	"github.com/ydmxcz/gds/stream"
)

func TestFlatMap(t *testing.T) {
	// sli := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	sli := slice.Of("11111", "22222", "33333", "44444", "55555", "66666", "77777", "88888", "99999")
	var count int64
	stream.Collect(stream.Inspect(
		stream.Parallel(stream.FlatMap(sli.Stream(),
			func(s string) stream.Stream[byte] {
				return slice.Of([]byte(s)...).Stream()
			}), 4),
		func(b byte) {
			fmt.Println(string(b))
		}),
		func(a byte) {
			atomic.AddInt64(&count, 1)
		})
	fmt.Println("count:", count)
}
