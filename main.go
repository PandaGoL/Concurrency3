package main

import (
	"fmt"
)

/*
Время попрактиковаться в работе в каналами.
Реализуем кольцевой буфер: записываем очередь в inCh, а потом читаем ее из outCh.
Нужно реализовать функционал кольцевого буфера, а именно функцию Run - при превышении лимитa (cap outCh), следует удалять старые значения очереди.
Пример: [0, 1, 2, 3, 4, 5, 6, 7, 8], размер буфера: 2, результат: [7, 8].
Из-за сложности изолировать задание в тестах пройдет так же решение: [6, 7, 8].
*/

func NewRingBuffer(inCh, outCh chan int) *ringBuffer {
	return &ringBuffer{
		inCh:  inCh,
		outCh: outCh,
	}
}

type ringBuffer struct {
	inCh  chan int
	outCh chan int
}

func (r *ringBuffer) Run() {
	// читаем из канала inCh.
	for v := range r.inCh {
		select {
		//пока пишется в канал outCh до блокировки, пишем
		case r.outCh <- v:
			//Если уже не лезет
		default:
			//выкидываем данные из канал
			<-r.outCh
			//пишем в освободившееся место
			r.outCh <- v
		}
	}
	// в конце закрываем outCh
	close(r.outCh)

}

func main() {
	max := 100
	inCh := make(chan int, max)
	outCh := make(chan int, 1)

	for i := 0; i < max; i++ {
		inCh <- i
	}

	rb := NewRingBuffer(inCh, outCh)
	close(inCh)
	rb.Run()

	resSlice := make([]int, 1)
	for res := range outCh {
		resSlice = append(resSlice, res)
	}
	fmt.Println(resSlice)
}
