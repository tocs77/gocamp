package main

import (
	"fmt"
)

func main() {
	var numbers = []int{1, 2, 3, 4, 5}
	var slice1 = numbers[0:3]
	slice2 := slice1

	slice1[0] = 10
	fmt.Println(slice1)
	fmt.Println(slice2)

}
