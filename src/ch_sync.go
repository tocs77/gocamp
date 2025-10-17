package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	numGoroutines := 3
	ch := make(chan int, numGoroutines)

	for i := range numGoroutines {
		go func(i int) {
			fmt.Printf("Goroutine %d started\n", i)
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			ch <- i
		}(i)
	}

	for range numGoroutines {
		fmt.Printf("Goroutine %d finished\n", <-ch)
	}

}
