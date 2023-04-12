package stream_test

import (
	"fmt"
	"testing"

	"github.com/ydmxcz/gds/collections/slice"
	"github.com/ydmxcz/gds/stream"
)

func TestMin(t *testing.T) {
	// sli := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	sli := slice.Of(61, 4, 6, 8, 60, 12, 14, 16, 18, 62)

	allMatch := stream.Min(sli.Stream(5))
	fmt.Println(allMatch)
}
