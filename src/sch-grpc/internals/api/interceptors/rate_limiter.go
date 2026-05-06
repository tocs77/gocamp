package interceptors

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
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

func clientIPFromAddr(addr net.Addr) string {
	if tcpAddr, ok := addr.(*net.TCPAddr); ok {
		return tcpAddr.IP.String()
	}
	host, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		return addr.String()
	}
	return host
}

func (rl *rateLimiter) RateLimiterInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get peer from context")
	}
	rl.mu.Lock()
	defer rl.mu.Unlock()
	clientIP := clientIPFromAddr(peer.Addr)
	count := rl.visitors[clientIP]
	rl.visitors[clientIP] = count + 1
	fmt.Printf("IP: %s, Count: %d\n", clientIP, rl.visitors[clientIP])
	if count >= rl.limit {
		return nil, status.Errorf(codes.ResourceExhausted, "Too many requests")
	}
	return handler(ctx, req)
}
