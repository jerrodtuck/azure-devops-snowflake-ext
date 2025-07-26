package middleware

import (
	"net/http"
	"time"
)

// RateLimitMiddleware implements rate limiting
func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	// Simple in-memory rate limiter (use Redis for production)
	clients := make(map[string][]time.Time)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientID := getClientIP(r)
			now := time.Now()

			// Clean old entries
			if requests, exists := clients[clientID]; exists {
				validRequests := []time.Time{}
				for _, t := range requests {
					if now.Sub(t) < time.Minute {
						validRequests = append(validRequests, t)
					}
				}
				clients[clientID] = validRequests
			}

			// Check rate limit
			if len(clients[clientID]) >= requestsPerMinute {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// Add current request
			clients[clientID] = append(clients[clientID], now)

			next.ServeHTTP(w, r)
		})
	}
}