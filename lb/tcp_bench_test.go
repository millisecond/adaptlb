package lb

import (
	"testing"
	"fmt"
)

func BenchmarkFib10(b *testing.B) {
	// run the Fib function b.N times
	i := 0
	for n := 0; n < b.N; n++ {
		i *= 2
	}
	fmt.Print(i)
}
