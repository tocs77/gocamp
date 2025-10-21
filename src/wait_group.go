package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//? =========== Basic example
// func main() {
// 	var wg sync.WaitGroup

// 	for i := range 3 {
// 		wg.Add(1)
// 		go worker(i, &wg)
// 	}

// 	wg.Wait()
// 	fmt.Println("All workers finished")
// }

// func worker(id int, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	fmt.Printf("Worker %d started\n", id)
// 	time.Sleep(time.Second)
// 	fmt.Printf("Worker %d finished\n", id)
// }

//? Example with channels

func worker(id int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d started\n", id)
	time.Sleep(time.Duration(rand.Intn(6)) * time.Second)
	fmt.Printf("Worker %d finished\n", id)
	results <- id * 2
}

func main() {
	var wg sync.WaitGroup
	results := make(chan int, 3)

	for i := range 3 {
		wg.Add(1)
		go worker(i, results, &wg)
	}
	go func() {
		wg.Wait()
		close(results)
		fmt.Println("All workers finished")
	}()

	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}
}
