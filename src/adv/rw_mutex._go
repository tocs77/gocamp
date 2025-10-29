package main

import (
	"fmt"
	"sync"
	"time"
)

var rwmu sync.RWMutex
var counter int

func readCounter(wg *sync.WaitGroup) {
	defer wg.Done()
	rwmu.RLock()
	time.Sleep(1 * time.Second)
	fmt.Println("Read counter: ", counter)
	rwmu.RUnlock()
}

func writeCouner(value int, wg *sync.WaitGroup) {
	defer wg.Done()
	rwmu.Lock()
	counter = value
	rwmu.Unlock()
}

func main() {
	var wg sync.WaitGroup

	for range 5 {
		wg.Add(1)
		go readCounter(&wg)
	}

	wg.Add(1)
	writeCouner(10, &wg)

	wg.Wait()
	fmt.Println("Final value in counter: ", counter)
}
