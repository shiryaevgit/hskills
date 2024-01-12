package main

import (
	"log"
	"time"
)

func main() {
	someCh := make(chan int)

	go func() { // пишем в канал
		defer close(someCh)
		for i := 0; i < 5; i++ {
			someCh <- i
		}
	}()

	for v := range BatchChan(someCh, time.Second, 3) {
		log.Println(v)
	}
}

func BatchChan(ch chan int, timeout time.Duration, limit int) chan []int {
	resch := make(chan []int)
	go func() {

		defer close(resch)
		batch := make([]int, 0, limit) // инициализацию батча с указанием capacity батча равным limit
		//timer := time.NewTimer(timeout)
		tiker := time.NewTicker(timeout)
		for {
			select {
			case val, ok := <-ch:
				if !ok {
					if len(batch) > 0 {
						resch <- batch
					}
					return
				}

				batch = append(batch, val)
				if len(batch) == limit {
					resch <- batch
					//batch=nil уточнить
					batch = make([]int, 0, limit)
					tiker.Reset(timeout)
				}

			case <-tiker.C:
				if len(batch) > 0 {
					resch <- batch
				}
				//return убираем, так как <-tiker.C, как ограничитель ожидания, а не сигнал для выхода из функции
			}
		}

	}()

	return resch
}
