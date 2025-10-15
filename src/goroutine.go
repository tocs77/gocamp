package main

import (
	"fmt"
	"time"
)

func main() {
	go sayHello()
	go printNumbers()
	go printLetters()
	fmt.Println("Hello, World! From main")
	var err error

	go func() {
		err = doBrokenWork()
	}()

	time.Sleep(2 * time.Second)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func sayHello() {
	time.Sleep(1 * time.Second)
	fmt.Println("Hello, World! From goroutine")
}

func printNumbers() {
	for i := range 10 {
		fmt.Println("Number: ", i)
		time.Sleep(100 * time.Millisecond)
	}
}

func printLetters() {
	for i := 'a'; i <= 'z'; i++ {
		fmt.Println("Letter: ", string(i))
		time.Sleep(200 * time.Millisecond)
	}
}

func doBrokenWork() error {
	return fmt.Errorf("broken work")
}
