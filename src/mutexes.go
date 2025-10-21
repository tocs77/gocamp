package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) getValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func mutexWworker(c *Counter, wg *sync.WaitGroup) {
	defer wg.Done()
	c.increment()
}

func noMutexWorker(c *Counter, wg *sync.WaitGroup) {
	defer wg.Done()
	c.value++
}

func main() {
	counter := Counter{}
	counter2 := Counter{}
	var wg sync.WaitGroup

	numWorkers := 10000

	for range numWorkers {
		wg.Add(2)
		go mutexWworker(&counter, &wg)
		go noMutexWorker(&counter2, &wg)
	}
	wg.Wait()
	fmt.Println("Final value in counter: ", counter.getValue())
	fmt.Println("Final value in counter2: ", counter2.getValue())
}
