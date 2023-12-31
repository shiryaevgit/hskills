package main

import (
	"fmt"
	"runtime"
)

func main() {
	nGoroutines := runtime.NumCPU()
	nElements := 1000000
	resSlices := createSlices(nElements, nGoroutines)

	fmt.Println(spreader(resSlices))
}

func calc(slice []int, ch chan int) {
	defer close(ch)

	sum := 0
	for _, v := range slice {
		sum += v
	}
	ch <- sum
}

func spreader(slices [][]int) int {
	sliceChan := make([]chan int, len(slices))
	for i := range sliceChan {
		sliceChan[i] = make(chan int)
	}

	for i := 0; i < len(slices); i++ {
		go calc(slices[i], sliceChan[i])
	}

	resSum := 0
	for _, ch := range sliceChan {
		resSum += <-ch
	}
	return resSum
}

func createSlices(nEl, nGor int) [][]int {
	slice := make([]int, nEl)
	for i := 0; i < len(slice); i++ {
		slice[i] = i + 1
	}

	slices := make([][]int, nGor)
	border := nEl / nGor

	for i := 0; i < nGor; i++ {
		if i == nGor-1 {
			slices[i] = slice
			break
		}

		slices[i] = slice[:border]
		slice = slice[border:]
	}
	return slices
}
