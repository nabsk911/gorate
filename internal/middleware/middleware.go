package middleware

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/nabsk911/gorate/internal/bucket"
)

func RateLimit(b *bucket.Bucket) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract IP from RemoteAddr (host:port format).
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)

			// Check if the request is allowed for the given IP.
			result, err := b.Allow(r.Context(), ip)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				next.ServeHTTP(w, r)
				return
			}

			// If rate limit exceeded, return 429 status.
			if !result.Allowed {
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{"error": "rate limit exceeded"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
