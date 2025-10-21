package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, job)
		time.Sleep(time.Second)
		fmt.Printf("Worker %d finished job %d\n", id, job)
		results <- job * 2
	}
}

func jobProducer(jobs chan<- int, jobsNumber int) {
	defer close(jobs)
	for j := 0; j < jobsNumber; j++ {
		jobs <- j
	}

}

func main() {
	workesNumber := 3
	jobsNumber := 10
	tasks := make(chan int, 3)
	results := make(chan int, 3)

	for w := range workesNumber {
		go worker(w, tasks, results)
	}

	go jobProducer(tasks, jobsNumber)

	for range jobsNumber {
		fmt.Printf("Result: %d\n", <-results)
	}

	close(results)
}
