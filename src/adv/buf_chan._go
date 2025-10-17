package main

import "fmt"

func main() {
	ch := make(chan int, 2)

	go func() {
		fmt.Println("Receiving from channel")
		fmt.Println(<-ch)
		fmt.Println(<-ch)
		fmt.Println(<-ch)
		fmt.Println("Done receiving from channel")
	}()

	ch <- 1
	ch <- 2
	ch <- 3

}
