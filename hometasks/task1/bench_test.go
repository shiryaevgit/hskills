package main

import (
	"runtime"
	"testing"
)

func generateData() [][]int {
	nGoroutines := runtime.NumCPU()
	nElements := 10000000
	return createSlices(nElements, nGoroutines)
}

func BenchmarkSpreader(b *testing.B) {
	slices := generateData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spreader(slices)
	}
}

func BenchmarkSpreader3(b *testing.B) {
	slices := generateData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spreader3(slices)
	}
}
