package main

import (
	"fmt"
	"sync"
	"time"
)

const bufferSize = 5

type Buffer struct {
	items []int
	mutex sync.Mutex
	cond  *sync.Cond
}

func NewBuffer(size int) *Buffer {
	b := &Buffer{
		items: make([]int, 0, size),
	}
	b.cond = sync.NewCond(&b.mutex)
	return b
}

func (b *Buffer) produce(item int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for len(b.items) == bufferSize {
		fmt.Println("Buffer is full, waiting for consumer ", item)
		b.cond.Wait()
	}
	b.items = append(b.items, item)
	fmt.Printf("Produced: %d\n", item)
	b.cond.Signal()
}

func (b *Buffer) consume() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for len(b.items) == 0 {
		b.cond.Wait()
	}
	item := b.items[0]
	b.items = b.items[1:]
	fmt.Printf("Consumed: %d\n", item)
	b.cond.Signal()
	return item
}

func producer(b *Buffer, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range 10 {
		b.produce(i + 100)
		time.Sleep(20 * time.Millisecond)
	}
}

func consumer(b *Buffer, wg *sync.WaitGroup) {
	defer wg.Done()
	for range 10 {
		b.consume()
		time.Sleep(600 * time.Millisecond)
	}
}

func main() {
	b := NewBuffer(bufferSize)
	var wg sync.WaitGroup
	wg.Add(2)
	go producer(b, &wg)
	go consumer(b, &wg)
	wg.Wait()
}
