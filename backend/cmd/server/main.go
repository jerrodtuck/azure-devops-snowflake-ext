package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"snowflake-dropdown-api/internal/api"
	"snowflake-dropdown-api/internal/config"
	"snowflake-dropdown-api/internal/database"
)

func main() {
	// Load .env file
	if err := config.LoadEnvFile(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Check required environment variables (skip if in TEST_MODE)
	if err := validateEnvironment(); err != nil {
		log.Fatalf("Environment validation failed: %v", err)
	}

	// Initialize database connection
	if err := database.InitializeDatabase(); err != nil {
		log.Printf("Warning: Database initialization failed: %v", err)
		log.Println("Server will start but database queries will fail")
		log.Println("Consider using TEST_MODE=true for testing without database")
	}

	// Load dynamic configuration
	if err := config.LoadConfig(); err != nil {
		log.Printf("Warning: Failed to load config: %v", err)
		log.Println("Using default configuration")
	}

	// Setup router and middleware
	handler := api.SetupRouter()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	logEndpoints()

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

// validateEnvironment checks required environment variables
func validateEnvironment() error {
	if os.Getenv("TEST_MODE") == "true" {
		log.Println("Running in TEST_MODE - using mock data")
		return nil
	}

	required := []string{
		"SNOWFLAKE_ACCOUNT",
		"SNOWFLAKE_USER",
		"SNOWFLAKE_DATABASE",
		"SNOWFLAKE_SCHEMA",
		"SNOWFLAKE_WAREHOUSE",
		"SNOWFLAKE_ROLE",
	}

	// Only require password if not using SSO
	if os.Getenv("SNOWFLAKE_AUTH_TYPE") != "externalbrowser" {
		required = append(required, "SNOWFLAKE_PASSWORD")
	}

	for _, env := range required {
		if os.Getenv(env) == "" {
			return fmt.Errorf("missing required environment variable: %s", env)
		}
	}

	return nil
}

// logEndpoints logs available API endpoints
func logEndpoints() {
	log.Printf("Available endpoints:")
	log.Printf("  GET /api/health - Health check")
	log.Printf("  GET /api/config - Get data types configuration")
	log.Printf("  GET /api/search/{type} - Search with dynamic data type")
	log.Printf("  GET /api/types - List available data types")
	log.Printf("  GET /api/dropdown/{type} - Legacy endpoint (redirects to search)")
	log.Printf("  POST /api/dynamic-search - Custom query endpoint")
	log.Printf("")
	log.Printf("Dynamic configuration loaded from: %s", getConfigFile())
}

// getConfigFile returns the configuration file path
func getConfigFile() string {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.json"
	}
	return configFile
}