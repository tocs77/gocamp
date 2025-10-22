package main

import (
	"fmt"
	"time"
)

type RateLimiter struct {
	tokens     chan struct{}
	refillTime time.Duration
}

func NewRateLimiter(rateLimit int, refillTime time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:     make(chan struct{}, rateLimit),
		refillTime: refillTime,
	}
	for range rateLimit {
		rl.tokens <- struct{}{}
	}
	go rl.startRefill()
	return rl
}

func (rl *RateLimiter) startRefill() {
	ticker := time.NewTicker(rl.refillTime)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case rl.tokens <- struct{}{}:
		default:
			// no-op
		}
	}
}

func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

func main() {
	rl := NewRateLimiter(5, 1*time.Second)
	for range 10 {
		if rl.Allow() {
			fmt.Println("Allowed")
		} else {
			fmt.Println("Not allowed")
		}
		time.Sleep(400 * time.Millisecond)
	}
}
