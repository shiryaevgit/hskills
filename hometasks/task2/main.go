package main

import (
	"log"
	"time"
)

func main() {
	someCh := make(chan int)

	go func() { // пишем в канал
		defer close(someCh)
		for i := 1; i < 50; i++ {
			someCh <- i
		}
	}()

	for v := range BatchChan(someCh, time.Second, 5) {
		log.Println(v)
	}
}

func BatchChan(ch chan int, timeout time.Duration, limit int) chan []int {
	resch := make(chan []int)
	go func() {
		// читать пачками
		// Limit - максимальный размер batch, по сколько считываем из канала за раз
		//Timeout - максимальное время ожидание непустого batch. Если не будет нового значения и таймер завершился, завершаем
		defer close(resch)
		batch := make([]int, 0)
		timer := time.NewTimer(timeout)

		for {
			select {
			case val, ok := <-ch: // если есть данные, то
				if !ok { // канал закрыт, проверяем есть ли данные, если есть добавляем в resch
					if len(batch) > 0 {
						resch <- batch
					}
					return // выходим
				}

				batch = append(batch, val) // если канал не закрыт, добавляем в batch
				if len(batch) == limit {   // проверка на лимит
					resch <- batch
					batch = nil          // обнулили
					timer.Reset(timeout) // таймер скинули
				}

			case <-timer.C: // если время ожидания вышло то, добрасываем в resch то что осталось и выходим
				if len(batch) > 0 {
					resch <- batch
				}
				return
			}
		}

	}()

	return resch
}
