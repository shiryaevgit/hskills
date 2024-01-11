package main

import (
	"reflect"
	"testing"
	"time"
)

func TestBatchChan(t *testing.T) {
	type args struct {
		ch      chan int
		timeout time.Duration
		limit   int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{name: "norm", args: args{ch: genChan(5), timeout: time.Second * 1, limit: 3}, want: [][]int{{0, 1, 2}, {3, 4}}},
		{name: "remainder 1", args: args{ch: genChan(7), timeout: time.Second * 1, limit: 3}, want: [][]int{{0, 1, 2}, {3, 4, 5}, {6}}},
		{name: "limit 1", args: args{ch: genChan(5), timeout: time.Second * 1, limit: 1}, want: [][]int{{0}, {1}, {2}, {3}, {4}}},
		{name: "limit nil", args: args{ch: genChan(5), timeout: time.Second * 1, limit: 0}, want: [][]int{{0, 1, 2, 3, 4}}},
		{name: "genChan(1)", args: args{ch: genChan(1), timeout: time.Second * 1, limit: 1}, want: [][]int{{0}}},
		{name: "genChan(1)", args: args{ch: genChan(0), timeout: time.Second * 1, limit: 1}, want: [][]int{}}, //  BatchChan()= [],want []
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BatchChan(tt.args.ch, tt.args.timeout, tt.args.limit)

			var res [][]int
			for batch := range got {
				res = append(res, batch)
			}

			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("BatchChan()= %v,want %v", res, tt.want)
			}

		})
	}
}

func genChan(n int) chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := 0; i < n; i++ {
			ch <- i
		}
	}()
	return ch
}
