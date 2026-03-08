package ratelimit

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type bucket struct {
	tokens   float64
	lastSeen time.Time
	mu       sync.Mutex
}

type Limiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     float64
	capacity float64
}

func NewLimiter(requestsPerMinute int) *Limiter {
	l := &Limiter{
		buckets:  make(map[string]*bucket),
		rate:     float64(requestsPerMinute) / 60.0,
		capacity: float64(requestsPerMinute),
	}
	go func() {
		for range time.Tick(time.Minute) {
			l.mu.Lock()
			for ip, b := range l.buckets {
				b.mu.Lock()
				if time.Since(b.lastSeen) > 5*time.Minute {
					delete(l.buckets, ip)
				}
				b.mu.Unlock()
			}
			l.mu.Unlock()
		}
	}()
	return l
}

func (l *Limiter) Allow(ip string) bool {
	l.mu.Lock()
	b, ok := l.buckets[ip]
	if !ok {
		b = &bucket{tokens: l.capacity, lastSeen: time.Now()}
		l.buckets[ip] = b
	}
	l.mu.Unlock()

	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastSeen).Seconds()
	b.tokens = min(l.capacity, b.tokens+elapsed*l.rate)
	b.lastSeen = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// Middleware returns an HTTP middleware that rate limits by IP.
func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realIP(r)
		if !l.Allow(ip) {
			http.Error(w, `{"success":false,"message":"rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuthRateLimit applies the given limiter only to auth login and register routes.
func AuthRateLimit(authLimiter *Limiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.Contains(path, "/auth/login") || strings.Contains(path, "/auth/register") {
			ip := realIP(r)
			if !authLimiter.Allow(ip + ":" + path) {
				http.Error(w, `{"success":false,"message":"muitas tentativas, aguarde"}`, http.StatusTooManyRequests)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func realIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		parts := strings.SplitN(fwd, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	if rip := r.Header.Get("X-Real-IP"); rip != "" {
		return strings.TrimSpace(rip)
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
