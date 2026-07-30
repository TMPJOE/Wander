package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/time/rate"
)

func TestShouldRateLimit(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"/api/v1/tours", true},
		{"/api/v1/uploads", true},
		{"/api/v1/health", true},
		{"/api/v1/", true},
		{"/api/", true},
		{"/assets/index-DfM-Q-Qg.js", false},
		{"/assets/ReviewCard-CJ6GTp7K.css", false},
		{"/uploads/abc123.jpg", false},
		{"/", false},
		{"/profile", false},  // SPA client route
		{"/tours/12", false}, // SPA client route (not /api/)
		{"", false},
		{"api/v1/tours", true}, // no leading slash tolerated
		{"/api", false},        // bare "/api" without trailing slash is not a limit-worthy prefix
		{"/api-docs", false},   // looks like api* but is not /api/
	}
	for _, c := range cases {
		if got := shouldRateLimit(c.path); got != c.want {
			t.Errorf("shouldRateLimit(%q) = %v, want %v", c.path, got, c.want)
		}
	}
}

func TestRateLimiterSkipsStaticAssets(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(1), 1) // tiny budget so it WOULD trip
	mw := RateLimiter(limiter)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	// Flood an asset path well beyond the burst budget; it must always reach
	// the handler, never return 429.
	for i := 0; i < 50; i++ {
		called = false
		req := httptest.NewRequest(http.MethodGet, "/assets/star-Wpf5clYo.js", nil)
		rec := httptest.NewRecorder()
		mw(next).ServeHTTP(rec, req)
		if !called {
			t.Fatalf("iteration %d: handler not called for asset path", i)
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("iteration %d: asset got %d, want 200", i, rec.Code)
		}
	}

	// An API path under the same flood MUST be throttled once the budget is
	// spent (demonstrating the limiter still guards the API).
	throttled := false
	for i := 0; i < 50; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tours", nil)
		rec := httptest.NewRecorder()
		mw(next).ServeHTTP(rec, req)
		if rec.Code == http.StatusTooManyRequests {
			throttled = true
			break
		}
	}
	if !throttled {
		t.Fatalf("API path was never throttled over 50 rapid requests; burst = 1")
	}
}
