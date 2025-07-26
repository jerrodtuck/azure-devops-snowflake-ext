package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "github.com/snowflakedb/gosnowflake"
)

// DropdownItem represents a single dropdown option
type DropdownItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// DropdownResponse represents the API response
type DropdownResponse struct {
	Data     []DropdownItem `json:"data"`
	Metadata Metadata       `json:"metadata"`
}

// Metadata contains information about the data
type Metadata struct {
	ExportedAt time.Time `json:"exported_at"`
	RowCount   int       `json:"row_count"`
	Source     string    `json:"source"`
	Cached     bool      `json:"cached"`
}

// Cache stores the dropdown data with expiration
type Cache struct {
	mu         sync.RWMutex
	data       map[string]DropdownResponse
	expiration time.Duration
}

// Global cache instance
var cache = &Cache{
	data:       make(map[string]DropdownResponse),
	expiration: 1 * time.Hour, // Cache for 1 hour
}

// Global database connection
var db *sql.DB

// Database connection string
func getConnectionString() string {
	account := os.Getenv("SNOWFLAKE_ACCOUNT")

	// Check if using SSO authentication
	if os.Getenv("SNOWFLAKE_AUTH_TYPE") == "externalbrowser" {
		// For SSO/External Browser auth, no password needed
		dsn := fmt.Sprintf("%s@%s/%s/%s?warehouse=%s&authenticator=externalbrowser",
			os.Getenv("SNOWFLAKE_USER"),
			account,
			os.Getenv("SNOWFLAKE_DATABASE"),
			os.Getenv("SNOWFLAKE_SCHEMA"),
			os.Getenv("SNOWFLAKE_WAREHOUSE"),
		)

		if role := os.Getenv("SNOWFLAKE_ROLE"); role != "" {
			dsn += "&role=" + role
		}

		log.Printf("Connecting to Snowflake using SSO (external browser) for account: %s", account)
		return dsn
	}

	// Standard username/password authentication
	password := url.QueryEscape(os.Getenv("SNOWFLAKE_PASSWORD"))

	dsn := fmt.Sprintf("%s:%s@%s/%s/%s?warehouse=%s",
		os.Getenv("SNOWFLAKE_USER"),
		password,
		account,
		os.Getenv("SNOWFLAKE_DATABASE"),
		os.Getenv("SNOWFLAKE_SCHEMA"),
		os.Getenv("SNOWFLAKE_WAREHOUSE"),
	)

	if role := os.Getenv("SNOWFLAKE_ROLE"); role != "" {
		dsn += "&role=" + role
	}

	log.Printf("Connecting to Snowflake account: %s", account)

	return dsn
}

// Get or set cache
func (c *Cache) Get(key string) (DropdownResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, exists := c.data[key]
	return val, exists
}

func (c *Cache) Set(key string, value DropdownResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value

	// Simple cache cleanup after expiration
	go func() {
		time.Sleep(c.expiration)
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
	}()
}

// Configuration will be loaded from config.json or environment

// DEPRECATED: Use handleSearch from endpoints.go instead
// This function is kept for backward compatibility only

// handleDropdownData redirects to the new dynamic search handler
func handleDropdownData(w http.ResponseWriter, r *http.Request) {
	// Use the new dynamic search handler for backward compatibility
	handleSearch(w, r)
}

// handleSearchData redirects to the new dynamic search handler
func handleSearchData(w http.ResponseWriter, r *http.Request) {
	// Use the new dynamic search handler for backward compatibility
	handleSearch(w, r)
}

// Search queries moved to config.json for dynamic configuration

// DEPRECATED: Query execution moved to endpoints.go with dynamic configuration

// handleHealth returns server health status
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// handleDataTypes is deprecated - use handleGetDataTypes from endpoints.go
func handleDataTypes(w http.ResponseWriter, r *http.Request) {
	// Redirect to the new dynamic handler
	handleGetDataTypes(w, r)
}

// getMockData returns mock data for testing
func getMockData(dataType string) (DropdownResponse, error) {
	var items []DropdownItem

	switch dataType {
	case "cc", "cost_centers":
		items = []DropdownItem{
			{Value: "1000", Label: "1000 - IT Department"},
			{Value: "2000", Label: "2000 - Finance Department"},
			{Value: "3000", Label: "3000 - Marketing Department"},
			{Value: "4000", Label: "4000 - Operations"},
			{Value: "5000", Label: "5000 - Human Resources"},
		}
	case "wbs":
		items = []DropdownItem{
			{Value: "WBS001", Label: "WBS001 - Project Alpha"},
			{Value: "WBS002", Label: "WBS002 - Project Beta"},
			{Value: "WBS003", Label: "WBS003 - Project Gamma"},
			{Value: "WBS004", Label: "WBS004 - Project Delta"},
			{Value: "WBS005", Label: "WBS005 - Project Epsilon"},
		}
	default:
		return DropdownResponse{}, fmt.Errorf("unknown data type: %s", dataType)
	}

	return DropdownResponse{
		Data: items,
		Metadata: Metadata{
			ExportedAt: time.Now().UTC(),
			RowCount:   len(items),
			Source:     dataType + " (mock)",
			Cached:     false,
		},
	}, nil
}

// getMockSearchData returns filtered mock data for testing
func getMockSearchData(dataType, searchTerm string) (DropdownResponse, error) {
	// Get all mock data
	allData, err := getMockData(dataType)
	if err != nil {
		return DropdownResponse{}, err
	}

	// Filter by search term
	searchLower := strings.ToLower(searchTerm)
	var filtered []DropdownItem

	for _, item := range allData.Data {
		if strings.Contains(strings.ToLower(item.Value), searchLower) ||
			strings.Contains(strings.ToLower(item.Label), searchLower) {
			filtered = append(filtered, item)
		}
	}

	return DropdownResponse{
		Data: filtered,
		Metadata: Metadata{
			ExportedAt: time.Now().UTC(),
			RowCount:   len(filtered),
			Source:     dataType + " (mock search)",
			Cached:     false,
		},
	}, nil
}

// initializeDatabase establishes the Snowflake connection with retry logic
func initializeDatabase() error {
	if os.Getenv("TEST_MODE") == "true" {
		log.Println("TEST_MODE enabled - skipping database connection")
		return nil
	}

	log.Println("Initializing Snowflake connection...")
	log.Println("Note: If using SSO, your browser will open for authentication")

	maxRetries := 3
	retryDelay := 5 * time.Second
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		var err error
		db, err = sql.Open("snowflake", getConnectionString())
		if err != nil {
			return fmt.Errorf("failed to create database connection: %v", err)
		}

		// Test the connection
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err = db.PingContext(ctx)
		cancel()

		if err == nil {
			log.Println("Successfully connected to Snowflake")
			return nil
		}

		// Check if it's an invalid token error (390195)
		if strings.Contains(err.Error(), "390195") || strings.Contains(err.Error(), "invalid") && strings.Contains(err.Error(), "token") {
			log.Printf("Authentication failed (attempt %d/%d): Invalid ID Token detected", attempt, maxRetries)
			
			if attempt < maxRetries {
				log.Printf("Retrying authentication in %v...", retryDelay)
				
				// Close the current connection before retry
				if db != nil {
					db.Close()
					db = nil
				}
				
				time.Sleep(retryDelay)
				continue
			}
		}

		// For other errors, don't retry
		return fmt.Errorf("failed to connect to Snowflake: %v", err)
	}

	return fmt.Errorf("failed to connect to Snowflake after %d attempts: invalid ID token", maxRetries)
}

func main() {
	// Load .env file
	if err := loadEnvFile(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Check required environment variables (skip if in TEST_MODE)
	if os.Getenv("TEST_MODE") != "true" {
		required := []string{
			"SNOWFLAKE_ACCOUNT",
			"SNOWFLAKE_USER",
		}

		// Only require password if not using SSO
		if os.Getenv("SNOWFLAKE_AUTH_TYPE") != "externalbrowser" {
			required = append(required, "SNOWFLAKE_PASSWORD")
		}

		required = append(required,
			"SNOWFLAKE_DATABASE",
			"SNOWFLAKE_SCHEMA",
			"SNOWFLAKE_WAREHOUSE",
			"SNOWFLAKE_ROLE",
		)

		for _, env := range required {
			if os.Getenv(env) == "" {
				log.Fatalf("Missing required environment variable: %s", env)
			}
		}
	} else {
		log.Println("Running in TEST_MODE - using mock data")
	}

	// Initialize database connection
	if err := initializeDatabase(); err != nil {
		log.Printf("Warning: Database initialization failed: %v", err)
		log.Println("Server will start but database queries will fail")
		log.Println("Consider using TEST_MODE=true for testing without database")
	}

	// Load dynamic configuration
	if err := LoadConfig(); err != nil {
		log.Printf("Warning: Failed to load config: %v", err)
		log.Println("Using default configuration")
	}

	// Create router
	router := mux.NewRouter()

	// API routes with subrouter for better organization
	api := router.PathPrefix("/api").Subrouter()
	
	// Health check
	api.HandleFunc("/health", handleHealth).Methods("GET", "OPTIONS")
	
	// Dynamic configuration endpoints
	api.HandleFunc("/config", handleGetConfig).Methods("GET", "OPTIONS")
	api.HandleFunc("/search/{type}", handleSearch).Methods("GET", "OPTIONS")
	api.HandleFunc("/types", handleGetDataTypes).Methods("GET", "OPTIONS")
	
	// Legacy endpoints for backward compatibility
	api.HandleFunc("/dropdown/{type}", handleDropdownData).Methods("GET", "OPTIONS")
	api.HandleFunc("/dropdown", handleDropdownData).Methods("GET", "OPTIONS")

	// Apply simple API key authentication if enabled
	var handler http.Handler = router

	if apiKey := os.Getenv("API_KEY"); apiKey != "" {
		log.Printf("API Key authentication enabled")
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health check
			if r.URL.Path == "/api/health" {
				router.ServeHTTP(w, r)
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

			router.ServeHTTP(w, r)
		})
	}

	// Enable CORS for Azure DevOps
	corsOrigins := []string{
		"https://dev.azure.com",
		"https://*.visualstudio.com",
		"https://*.gallery.vsassets.io", // Azure DevOps extension gallery
		"https://*.gallerycdn.vsassets.io", // Azure DevOps extension gallery CDN
		"http://localhost:*", // For testing
	}

	// Allow custom CORS origins from environment
	if customOrigins := os.Getenv("CORS_ORIGINS"); customOrigins != "" {
		corsOrigins = strings.Split(customOrigins, ",")
	}

	c := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			// Allow any origin from Azure DevOps gallery
			if strings.Contains(origin, ".gallery.vsassets.io") || strings.Contains(origin, ".gallerycdn.vsassets.io") {
				return true
			}
			// Allow any origin from Visual Studio
			if strings.Contains(origin, ".visualstudio.com") {
				return true
			}
			// Allow dev.azure.com
			if origin == "https://dev.azure.com" {
				return true
			}
			// Allow localhost for testing
			if strings.HasPrefix(origin, "http://localhost:") {
				return true
			}
			// Check against custom origins
			for _, allowed := range corsOrigins {
				if origin == allowed {
					return true
				}
			}
			return false
		},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		AllowCredentials: false,
	})

	handler = c.Handler(handler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  GET /api/health - Health check")
	log.Printf("  GET /api/config - Get data types configuration")
	log.Printf("  GET /api/search/{type} - Search with dynamic data type")
	log.Printf("  GET /api/types - List available data types")
	log.Printf("  GET /api/dropdown/{type} - Legacy endpoint (redirects to search)")
	log.Printf("")
	log.Printf("Dynamic configuration loaded from: %s", os.Getenv("CONFIG_FILE"))

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
