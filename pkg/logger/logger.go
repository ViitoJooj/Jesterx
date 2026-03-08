package logger

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	maxLines = 500
	logFile  = "logs.txt"
)

var mu sync.Mutex

type captureWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (cw *captureWriter) WriteHeader(status int) {
	cw.status = status
	cw.ResponseWriter.WriteHeader(status)
}

func (cw *captureWriter) Write(b []byte) (int, error) {
	n, err := cw.ResponseWriter.Write(b)
	cw.size += n
	return n, err
}

type entry struct {
	ts           time.Time
	ip           string
	username     string
	method       string
	route        string
	status       int
	duration     time.Duration
	requestSize  int64
	responseSize int
	headers      http.Header
}

func write(e entry) {
	user := e.username
	if user == "" {
		user = "anon"
	}

	var hParts []string
	for k, v := range e.headers {
		lower := strings.ToLower(k)
		if lower == "authorization" || lower == "cookie" {
			continue
		}
		hParts = append(hParts, k+"="+strings.Join(v, ","))
	}

	line := fmt.Sprintf(
		"[%s] IP=%-15s User=%-36s %s %-40s Status=%d Time=%dms Sent=%dB Recv=%dB Headers=[%s]",
		e.ts.Format("2006-01-02 15:04:05"),
		e.ip,
		user,
		e.method,
		e.route,
		e.status,
		e.duration.Milliseconds(),
		e.requestSize,
		e.responseSize,
		strings.Join(hParts, " | "),
	)

	mu.Lock()
	defer mu.Unlock()

	var lines []string
	if f, err := os.Open(logFile); err == nil {
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			lines = append(lines, sc.Text())
		}
		f.Close()
	}

	lines = append(lines, line)
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}

	f, err := os.Create(logFile)
	if err != nil {
		return
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	for _, l := range lines {
		fmt.Fprintln(bw, l)
	}
	bw.Flush()
}

func getIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.SplitN(xff, ",", 2)[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

type UserFromContext func(ctx context.Context) string

func Middleware(getUser UserFromContext) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			var reqSize int64
			if r.Body != nil && r.Body != http.NoBody {
				buf, err := io.ReadAll(r.Body)
				if err == nil {
					reqSize = int64(len(buf))
					r.Body = io.NopCloser(bytes.NewReader(buf))
				}
			}

			cw := &captureWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(cw, r)

			var username string
			if getUser != nil {
				username = getUser(r.Context())
			}

			go write(entry{
				ts:           start,
				ip:           getIP(r),
				username:     username,
				method:       r.Method,
				route:        r.URL.Path,
				status:       cw.status,
				duration:     time.Since(start),
				requestSize:  reqSize,
				responseSize: cw.size,
				headers:      r.Header.Clone(),
			})
		})
	}
}
