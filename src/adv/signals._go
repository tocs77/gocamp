package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	pid := os.Getpid()
	fmt.Println("PID: ", pid)
	sigs := make(chan os.Signal, 1)

	//Notify the channel on interrupt and terminate signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		sig := <-sigs
		switch sig {
		case syscall.SIGINT:
			fmt.Println("Received interrupt signal")
		case syscall.SIGTERM:
			fmt.Println("Received terminate signal")
		case syscall.SIGHUP:
			fmt.Println("Received hang up signal")
		}
	}()

	//simulate long running process
	fmt.Println("Running long running process...")
	for {
		time.Sleep(1 * time.Second)
	}
}
