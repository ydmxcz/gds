package stream_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ydmxcz/gds/collections/slice"
	"github.com/ydmxcz/gds/stream"
)

func TestFunc(t *testing.T) {
	sli := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1666, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1888, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1999, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	s := sli.Stream(8)

	stream.Collect(stream.Filter(s, func(a int) bool {
		time.Sleep(300 * time.Millisecond)
		return a%2 == 0
	}), func(a int) {
		fmt.Println(a)
	})
}
