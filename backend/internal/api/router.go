package api

import (
	"log"
	"net/http"
	"os"

	"snowflake-dropdown-api/internal/api/handlers"
	"snowflake-dropdown-api/internal/api/middleware"
	"snowflake-dropdown-api/internal/config"

	"github.com/gorilla/mux"
)

// SetupRouter configures and returns the HTTP router
func SetupRouter() http.Handler {
	router := mux.NewRouter()

	// API routes with subrouter for better organization
	api := router.PathPrefix("/api").Subrouter()

	// Health check
	api.HandleFunc("/health", handlers.HandleHealth).Methods("GET", "OPTIONS")

	// Dynamic configuration endpoints
	api.HandleFunc("/config", handlers.HandleGetConfig).Methods("GET", "OPTIONS")
	api.HandleFunc("/search/{type}", handlers.HandleSearch).Methods("GET", "OPTIONS")
	api.HandleFunc("/types", handlers.HandleGetDataTypes).Methods("GET", "OPTIONS")

	/*     // Legacy endpoints for backward compatibility
	       api.HandleFunc("/dropdown/{type}", handlers.HandleDropdownData).Methods("GET", "OPTIONS")
	       api.HandleFunc("/dropdown", handlers.HandleDropdownData).Methods("GET", "OPTIONS") */

	// Dynamic search endpoint (POST for custom queries)
	api.HandleFunc("/dynamic-search", handlers.HandleDynamicSearch).Methods("POST", "OPTIONS")

	// Apply middleware stack
	var handler http.Handler = router

	// Apply CORS middleware first to handle preflight requests
	corsHandler := middleware.SetupCORS()
	handler = corsHandler.Handler(handler)

	// Apply authentication middleware if configured
	if shouldUseSimpleAuth() {
		log.Printf("Simple API Key authentication enabled")
		handler = middleware.SimpleAPIKeyMiddleware()(handler)
	} else if shouldUseAdvancedAuth() {
		log.Printf("Advanced authentication enabled")
		secConfig := config.LoadSecurityConfig()
		handler = middleware.AuthMiddleware(secConfig)(handler)
	}

	// Apply rate limiting if configured
	if rateLimit := getRateLimit(); rateLimit > 0 {
		log.Printf("Rate limiting enabled: %d requests per minute", rateLimit)
		handler = middleware.RateLimitMiddleware(rateLimit)(handler)
	}

	return handler
}

// shouldUseSimpleAuth checks if simple API key auth should be used
func shouldUseSimpleAuth() bool {
	return os.Getenv("API_KEY") != ""
}

// shouldUseAdvancedAuth checks if advanced authentication should be used
func shouldUseAdvancedAuth() bool {
	return os.Getenv("AUTH_ENABLED") == "true" ||
		os.Getenv("JWT_ENABLED") == "true" ||
		os.Getenv("IP_WHITELIST") != ""
}

// getRateLimit returns the configured rate limit (0 = disabled)
func getRateLimit() int {
	// Could be made configurable via environment variable
	return 0 // Disabled by default
}

// LogRoutes logs all registered routes for debugging
func LogRoutes(router *mux.Router) {
	log.Printf("Available endpoints:")
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			methods, _ := route.GetMethods()
			if len(methods) > 0 {
				log.Printf("  %v %s", methods, pathTemplate)
			}
		}
		return nil
	})
}
