package main

import "fmt"

func main() {
	greeting := make(chan string)

	greetStringMessages := []string{"Hello, World!", "Hello, Go!", "Hello, Channels!"}

	go func() {
		for _, greetString := range greetStringMessages {
			greeting <- greetString
		}
		close(greeting)
	}()

	for greetString := range greeting {
		fmt.Println(greetString)
	}

}
