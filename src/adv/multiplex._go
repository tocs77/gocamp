package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- 1
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- 2
	}()

	select {
	case v := <-ch1:
		fmt.Println("Received from ch1: ", v)
	case v := <-ch2:
		fmt.Println("Received from ch2: ", v)
		// default:
		// 	fmt.Println("No channel received")
	}
}
