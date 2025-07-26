package main

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SecurityConfig holds security settings
type SecurityConfig struct {
	APIKeyEnabled bool
	APIKeys       []string
	JWTEnabled    bool
	JWTSecret     string
	IPWhitelist   []string
}

// LoadSecurityConfig loads security settings from environment
func LoadSecurityConfig() SecurityConfig {
	config := SecurityConfig{
		APIKeyEnabled: os.Getenv("AUTH_ENABLED") == "true",
		JWTEnabled:    os.Getenv("JWT_ENABLED") == "true",
		JWTSecret:     os.Getenv("JWT_SECRET"),
	}

	// Load API keys
	if keys := os.Getenv("API_KEYS"); keys != "" {
		config.APIKeys = strings.Split(keys, ",")
	}

	// Load IP whitelist
	if ips := os.Getenv("IP_WHITELIST"); ips != "" {
		config.IPWhitelist = strings.Split(ips, ",")
	}

	return config
}

// AuthMiddleware handles authentication
func AuthMiddleware(config SecurityConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health checks
			if r.URL.Path == "/api/health" {
				next.ServeHTTP(w, r)
				return
			}

			// IP Whitelist check
			if len(config.IPWhitelist) > 0 {
				clientIP := getClientIP(r)
				if !isIPAllowed(clientIP, config.IPWhitelist) {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}

			// API Key authentication
			if config.APIKeyEnabled {
				apiKey := r.Header.Get("X-API-Key")
				if apiKey == "" {
					apiKey = r.URL.Query().Get("apikey")
				}

				if !isValidAPIKey(apiKey, config.APIKeys) {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}

			// JWT authentication
			if config.JWTEnabled {
				tokenString := extractToken(r)
				if tokenString == "" {
					http.Error(w, "Missing token", http.StatusUnauthorized)
					return
				}

				if !isValidJWT(tokenString, config.JWTSecret) {
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

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

// Helper functions

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return strings.Split(r.RemoteAddr, ":")[0]
}

func isIPAllowed(ip string, whitelist []string) bool {
	for _, allowed := range whitelist {
		if ip == allowed {
			return true
		}
	}
	return false
}

func isValidAPIKey(key string, validKeys []string) bool {
	for _, validKey := range validKeys {
		if subtle.ConstantTimeCompare([]byte(key), []byte(validKey)) == 1 {
			return true
		}
	}
	return false
}

func extractToken(r *http.Request) string {
	// Check Authorization header
	bearerToken := r.Header.Get("Authorization")
	if strings.HasPrefix(bearerToken, "Bearer ") {
		return strings.TrimPrefix(bearerToken, "Bearer ")
	}

	// Check query parameter
	return r.URL.Query().Get("token")
}

func isValidJWT(tokenString, secret string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	return err == nil && token.Valid
}

// GenerateAPIKey generates a secure API key
func GenerateAPIKey() string {
	// In production, use crypto/rand
	return "sk_live_" + generateRandomString(32)
}

func generateRandomString(length int) string {
	// Implementation omitted for brevity
	// Use crypto/rand in production
	return "example_random_string"
}
