package safeguard

import (
	"log"
	"net"
	"net/http"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

func realIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return strings.TrimSpace(strings.SplitN(fwd, ",", 2)[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

type banEntry struct {
	until   time.Time
	strikes int
}

type IPBanner struct {
	mu      sync.Mutex
	bans    map[string]*banEntry
	maxHits int
	banFor  time.Duration
}

func NewIPBanner(maxHits int, banDuration time.Duration) *IPBanner {
	b := &IPBanner{
		bans:    make(map[string]*banEntry),
		maxHits: maxHits,
		banFor:  banDuration,
	}
	go func() {
		for range time.Tick(5 * time.Minute) {
			now := time.Now()
			b.mu.Lock()
			for ip, entry := range b.bans {
				if !entry.until.IsZero() && now.After(entry.until) {
					delete(b.bans, ip)
				}
			}
			b.mu.Unlock()
		}
	}()
	return b
}

func (b *IPBanner) Strike(ip string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	entry, ok := b.bans[ip]
	if !ok {
		entry = &banEntry{}
		b.bans[ip] = entry
	}
	entry.strikes++
	if entry.strikes >= b.maxHits {
		entry.until = time.Now().Add(b.banFor)
		entry.strikes = 0
	}
}

func (b *IPBanner) IsBanned(ip string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	entry, ok := b.bans[ip]
	if !ok {
		return false
	}
	return !entry.until.IsZero() && time.Now().Before(entry.until)
}

func (b *IPBanner) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if b.IsBanned(realIP(r)) {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"success":false,"message":"acesso temporariamente bloqueado"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[PANIC] %v\n%s", rec, debug.Stack())
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, `{"success":false,"message":"internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func BodyLimit(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxBytes {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, `{"success":false,"message":"payload too large"}`, http.StatusRequestEntityTooLarge)
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}

func PathTraversalGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawPath := r.URL.Path
		if strings.Contains(rawPath, "..") {
			http.Error(w, `{"success":false,"message":"invalid path"}`, http.StatusBadRequest)
			return
		}
		cleaned := path.Clean("/" + rawPath)
		if cleaned != rawPath && cleaned+"/" != rawPath {
			http.Error(w, `{"success":false,"message":"invalid path"}`, http.StatusBadRequest)
			return
		}
		if strings.Contains(r.URL.RawQuery, "..") {
			http.Error(w, `{"success":false,"message":"invalid request"}`, http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func PaginationGuard(maxLimit int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			modified := false
			for _, key := range []string{"limit", "per_page", "page_size"} {
				v := q.Get(key)
				if v == "" {
					continue
				}
				n, err := strconv.Atoi(v)
				if err != nil || n < 1 {
					q.Set(key, "1")
					modified = true
				} else if n > maxLimit {
					q.Set(key, strconv.Itoa(maxLimit))
					modified = true
				}
			}
			if modified {
				r.URL.RawQuery = q.Encode()
			}
			next.ServeHTTP(w, r)
		})
	}
}
