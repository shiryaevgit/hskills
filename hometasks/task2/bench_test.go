package main

import (
	"testing"
	"time"
)

func BenchmarkBatchChan(b *testing.B) {
	timeout := time.Second
	limit := 5
	someCh := generateData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		BatchChan(someCh, timeout, limit)
	}
}

func generateData() chan int {
	someCh := make(chan int)

	go func() {
		defer close(someCh)
		for i := 1; i < 50; i++ {
			someCh <- i
		}
	}()
	return someCh
}
