package stream_test

import (
	"fmt"
	"testing"

	"github.com/ydmxcz/gds/collections/slice"
	"github.com/ydmxcz/gds/stream"
)

func TestAllMatch(t *testing.T) {
	// sli := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	sli := slice.Of(2, 4, 6, 8, 10, 12, 14, 16, 18, 21)

	allMatch := stream.AllMatch(sli.Stream(), func(n int) bool {
		return n%2 == 0
	})
	fmt.Println(allMatch)
}

func TestAnyMatch(t *testing.T) {
	sli := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	// sli := slice.Of(1, 3, 5, 7, 9, 11, 13, 15, 17, 19)
	// sli := slice.Of(2, 4, 6, 8, 10, 12, 14, 16, 18, 21)

	allMatch := stream.AnyMatch(sli.Stream(), func(n int) bool {
		return n%2 == 0
	})
	fmt.Println(allMatch)
}
