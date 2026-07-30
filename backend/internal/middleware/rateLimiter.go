package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
	"wander/backend/internal/utils"

	"golang.org/x/time/rate"
)

// Each visitors has its own rate limiter
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IPRateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		visitors: make(map[string]*visitor),
		rate:     r,
		burst:    b,
	}
	go limiter.cleanupVisitors()

	return limiter
}

// Add visitor rate limiter
func (i *IPRateLimiter) getVisitor(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	v, exists := i.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(i.rate, i.burst)
		i.visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

// Every minute delete visitors that have not been seen in 5 minutes
func (i *IPRateLimiter) cleanupVisitors() {
	for {
		time.Sleep(1 * time.Minute)

		i.mu.Lock()
		for ip, v := range i.visitors {
			if time.Since(v.lastSeen) > 5*time.Minute {
				delete(i.visitors, ip)
			}
		}
		i.mu.Unlock()
	}
}

func RateLimiter(ipLimiter *IPRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only protect the API. Static assets (/assets/*, /uploads/*) and
			// the SPA fallback ("/*" -> index.html) must never be rate limited:
			// a cold page load fetches dozens of hashed JS/CSS files in
			// parallel from one IP and would otherwise blow the burst budget,
			// causing 429s whose JSON body browsers then reject as
			// "NS_ERROR_CORRUPTED_CONTENT" / "disallowed MIME type".
			if !shouldRateLimit(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Extract IP address from visitor
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			// Get rate limiter of this visitor
			limiter := ipLimiter.getVisitor(ip)

			if !limiter.Allow() {
				utils.SendError(w, http.StatusTooManyRequests, "Too many Request", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// shouldRateLimit reports whether a given request path should be subject to
// the per-IP rate limiter. Vite emits hashed module files under /assets/* and
// the SPA serves index.html at the root, so everything outside /api/ is a
// static file or client-side route, not a backend workload to throttle.
func shouldRateLimit(path string) bool {
	// The API is versioned under /api/v1/. Anything else is served by the
	// FileServer or the SPA fallback and bypasses the limiter. We compare on
	// the trimmed path so "/api/" and "/api" behave the same. Keeping this on
	// the literal "/api/" prefix (not "/api/v1/") future-proofs new minor
	// versions without reopening the limiter.
	p := path
	if len(p) > 0 && p[0] == '/' {
		p = p[1:]
	}
	return len(p) >= 4 && p[:4] == "api/"
}
