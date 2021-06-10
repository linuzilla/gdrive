package test_code

import (
	"fmt"
	"testing"
	"time"
)

func TestElapse(t *testing.T) {
	started := time.Now()
	last := time.Now()
	var lastRead int64
	var currentRead int64

	for i := 0; ; i += 1 {

		currentRead += 10000
		time.Sleep(time.Millisecond * 100)
		if i%150 == 149 {
			now := time.Now()

			elapsedTotal := now.Sub(started) / time.Second
			elapsedLast := now.Sub(last) / time.Second
			periodRead := currentRead - lastRead

			avgRead := float64(currentRead) / float64(elapsedTotal)
			currRead := float64(periodRead) / float64(elapsedLast)

			fmt.Printf("reading %d, (current %.2f, ) (avg: %.2f) (bytes/sec) %d, %d\n",
				currentRead,
				currRead, avgRead,
				elapsedTotal, elapsedLast)

			lastRead = currentRead
			last = now
		}
	}
}
