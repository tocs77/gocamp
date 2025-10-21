package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type AtomicCounter struct {
	count int64
}

func (ac *AtomicCounter) Increment(wg *sync.WaitGroup) {
	defer wg.Done()
	atomic.AddInt64(&ac.count, 1)
}

func (ac *AtomicCounter) GetValue() int64 {
	return atomic.LoadInt64(&ac.count)
}

func main() {
	counter := AtomicCounter{}
	var wg sync.WaitGroup

	numWorkers := 10000

	for range numWorkers {
		wg.Add(1)
		go counter.Increment(&wg)
	}
	wg.Wait()
	fmt.Println("Final value in counter: ", counter.GetValue())
}
