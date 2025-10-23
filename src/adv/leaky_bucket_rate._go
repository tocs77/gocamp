package main

import (
	"fmt"
	"sync"
	"time"
)

type LeakyBucketRateLimiter struct {
	capacity int
	leakRate time.Duration
	tokens   int
	lastLeak time.Time
	mu       sync.Mutex
}

func NewLeakyBucketRateLimiter(capacity int, leakRate time.Duration) *LeakyBucketRateLimiter {
	return &LeakyBucketRateLimiter{
		capacity: capacity,
		leakRate: leakRate,
		tokens:   capacity,
		lastLeak: time.Now(),
	}
}

func (lbrl *LeakyBucketRateLimiter) Allow() bool {
	lbrl.mu.Lock()
	defer lbrl.mu.Unlock()

	now := time.Now()
	elapsedTime := now.Sub(lbrl.lastLeak)
	tokensToAdd := int(elapsedTime.Seconds() / lbrl.leakRate.Seconds())

	// Store original state for printing
	originalTokens := lbrl.tokens

	lbrl.tokens = min(lbrl.tokens+tokensToAdd, lbrl.capacity)
	lbrl.lastLeak = lbrl.lastLeak.Add(time.Duration(tokensToAdd) * lbrl.leakRate)

	// Create visual bucket representation
	bucketDisplay := ""
	for i := 0; i < lbrl.capacity; i++ {
		if i < lbrl.tokens {
			bucketDisplay += "●"
		} else {
			bucketDisplay += "○"
		}
	}

	fmt.Printf("┌─── Leaky Bucket State ──┐\n")
	fmt.Printf("│ Elapsed: %8.1fms     │\n", elapsedTime.Seconds()*1000)
	fmt.Printf("│ Before:  %2d tokens      │\n", originalTokens)
	fmt.Printf("│ Added:   %2d tokens      │\n", tokensToAdd)
	fmt.Printf("│ Current: %2d/%d tokens    │\n", lbrl.tokens, lbrl.capacity)
	fmt.Printf("│ Bucket:  [%s]        │\n", bucketDisplay)

	if lbrl.tokens < 1 {
		fmt.Printf("│ Result:  ❌ DENIED      │\n")
		fmt.Printf("└─────────────────────────┘\n\n")
		return false
	}

	lbrl.tokens--
	fmt.Printf("│ Consumed: 1 token       │\n")
	fmt.Printf("│ Result:  ✅ ALLOWED     │\n")
	fmt.Printf("└─────────────────────────┘\n\n")
	return true
}

func main() {
	fmt.Println("🪣 Leaky Bucket Rate Limiter Demo")
	fmt.Println("Capacity: 5 tokens, Leak Rate: 500ms per token")
	fmt.Println("Requesting every 200ms...")

	lbrl := NewLeakyBucketRateLimiter(5, 500*time.Millisecond)
	for i := range 15 {
		fmt.Printf("Request #%d:\n", i+1)
		lbrl.Allow()
		time.Sleep(200 * time.Millisecond)
	}
}
