package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {

	bgContext := context.Background()
	mainContext, mainCancel := context.WithTimeout(bgContext, 1*time.Second)
	defer mainCancel()
	ctx, cancel := context.WithCancel(mainContext)
	ch := make(chan string)
	go doWork(ctx, ch)
	go switchContext(cancel)
	for v := range ch {
		fmt.Println(v)
	}

}

func doWork(ctx context.Context, ch chan<- string) {
	defer close(ch)
	counter := 0
	for {
		select {
		case <-ctx.Done():
			return
		case ch <- fmt.Sprintf("Work %d", counter):
			counter++
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func switchContext(cancel context.CancelFunc) {
	time.Sleep(time.Duration(rand.Intn(6)) * time.Second)
	cancel()
	fmt.Println("Switching context")
}
