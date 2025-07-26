package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"snowflake-dropdown-api/internal/config"
)

// AuthMiddleware handles authentication
func AuthMiddleware(secConfig config.SecurityConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health checks
			if r.URL.Path == "/api/health" {
				next.ServeHTTP(w, r)
				return
			}

			// IP Whitelist check
			if len(secConfig.IPWhitelist) > 0 {
				clientIP := getClientIP(r)
				if !isIPAllowed(clientIP, secConfig.IPWhitelist) {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}

			// API Key authentication
			if secConfig.APIKeyEnabled {
				apiKey := r.Header.Get("X-API-Key")
				if apiKey == "" {
					apiKey = r.URL.Query().Get("apikey")
				}

				if !isValidAPIKey(apiKey, secConfig.APIKeys) {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}

			// JWT authentication
			if secConfig.JWTEnabled {
				tokenString := extractToken(r)
				if tokenString == "" {
					http.Error(w, "Missing token", http.StatusUnauthorized)
					return
				}

				if !isValidJWT(tokenString, secConfig.JWTSecret) {
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// SimpleAPIKeyMiddleware provides basic API key authentication (legacy support)
func SimpleAPIKeyMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health check
			if r.URL.Path == "/api/health" {
				next.ServeHTTP(w, r)
				return
			}

			apiKey := os.Getenv("API_KEY")
			if apiKey == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Check API key
			providedKey := r.Header.Get("X-API-Key")
			if providedKey == "" {
				providedKey = r.URL.Query().Get("apikey")
			}

			if providedKey != apiKey {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

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