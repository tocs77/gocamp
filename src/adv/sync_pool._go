package main

import (
	"fmt"
	"sync"
)

type Person struct {
	Name string
	Age  int
}

func (p *Person) String() string {
	return fmt.Sprintf("Name: %s, Age: %d", p.Name, p.Age)
}

func main() {
	var pool = sync.Pool{
		New: func() any {
			fmt.Println("Creating new person")
			return &Person{}
		},
	}

	//get from pool
	person1 := pool.Get().(*Person)
	person1.Name = "Jane"
	person1.Age = 20
	fmt.Println("Person1: ", person1)

	pool.Put(person1)
	fmt.Println("Person1 put back to pool")

	person2 := pool.Get().(*Person)
	fmt.Println("Got from pool person2: ", person2)

	person3 := pool.Get().(*Person)
	fmt.Println("Got from pool person3: ", person3)

	pool.Put(person2)
	pool.Put(person3)
	fmt.Println("Person2 and person3 put back to pool")

	person4 := pool.Get().(*Person)
	fmt.Println("Got from pool person4: ", person4)
}
