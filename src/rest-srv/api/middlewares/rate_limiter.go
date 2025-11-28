package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mu        sync.Mutex
	visitors  map[string]int
	limit     int
	resetTime time.Duration
}

func NewRateLimiter(limit int, resetTime time.Duration) *rateLimiter {
	rl := &rateLimiter{
		limit:     limit,
		resetTime: resetTime,
		visitors:  make(map[string]int),
	}
	go rl.resetVisitorCount()
	return rl
}

func (rl *rateLimiter) resetVisitorCount() {
	for {
		time.Sleep(rl.resetTime)
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract IP address from RemoteAddr (removes port)
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// If SplitHostPort fails, use RemoteAddr as-is (fallback)
			ip = r.RemoteAddr
		}
		rl.mu.Lock()
		defer rl.mu.Unlock()
		count := rl.visitors[ip]
		rl.visitors[ip] = count + 1
		fmt.Printf("IP: %s, Count: %d\n", ip, rl.visitors[ip])
		if count >= rl.limit {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
