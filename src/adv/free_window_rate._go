package main

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	mu        sync.Mutex
	count     int
	limit     int
	window    time.Duration
	resetTime time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	if now.After(rl.resetTime) {
		rl.count = 0
		rl.resetTime = now.Add(rl.window)
	}
	if rl.count >= rl.limit {
		return false
	}
	rl.count++
	return true
}

func main() {
	rl := NewRateLimiter(5, 2*time.Second)
	for range 20 {
		if rl.Allow() {
			fmt.Println("Allowed")
		} else {
			fmt.Println("Not allowed")
		}
		time.Sleep(200 * time.Millisecond)
	}
}
