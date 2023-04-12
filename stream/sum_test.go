package stream_test

import (
	"fmt"
	"testing"

	"github.com/ydmxcz/gds/collections/slice"
	"github.com/ydmxcz/gds/stream"
)

func TestSum(t *testing.T) {
	sli := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	sum := stream.Sum(sli.Stream())
	fmt.Println(sum)
}

func TestSum_Parallel(t *testing.T) {
	sli := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	sum := stream.Sum(stream.Parallel(sli.Stream(), 8))
	fmt.Println(sum)
}
