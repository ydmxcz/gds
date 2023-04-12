package stream_test

import (
	"fmt"
	"testing"

	"github.com/ydmxcz/gds/collections/slice"
	"github.com/ydmxcz/gds/stream"
)

func TestSkip(t *testing.T) {
	stm := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Stream(4)

	stream.Collect(stream.Skip(stm, 3), func(a int) {
		fmt.Print(a, " ")
	})

}

func TestSkipWhile(t *testing.T) {
	stm := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Stream(4)

	stream.Collect(stream.SkipWhile(stm, func(a int) bool {
		return a%2 == 0
	}), func(a int) {
		fmt.Print(a, " ")
	})

}
