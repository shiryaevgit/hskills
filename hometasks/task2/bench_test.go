package main

import (
	"testing"
	"time"
)

func BenchmarkBatchChan(b *testing.B) {
	b.ReportAllocs()

	timeout := time.Second
	limit := 5
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		someCh := generateData()
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
