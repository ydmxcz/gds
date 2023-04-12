package stream_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ydmxcz/gds/collections/slice"
	"github.com/ydmxcz/gds/stream"
)

func TestWorldCount(t *testing.T) {
	s := slice.Of("hello world", "hello world mcz", "hello world", "hello world", "hello world").Stream()
	count := stream.ReduceWith(stream.Map(s,
		func(str string) []string {
			return strings.Split(str, " ")
		}), 0,
		func(count int, val []string) int {
			return count + len(val)
		},
		func(accum, val int) int {
			return accum + val
		})
	fmt.Println(count)
}
