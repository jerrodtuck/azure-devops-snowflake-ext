package middleware

import (
	"os"
	"strings"

	"github.com/rs/cors"
)

// SetupCORS configures CORS middleware for Azure DevOps
func SetupCORS() *cors.Cors {
	corsOrigins := []string{
		"https://dev.azure.com",
		"https://*.visualstudio.com",
		"https://*.gallery.vsassets.io",    // Azure DevOps extension gallery
		"https://*.gallerycdn.vsassets.io", // Azure DevOps extension gallery CDN
		"http://localhost:*",               // For testing
	}

	// Allow custom CORS origins from environment
	if customOrigins := os.Getenv("CORS_ORIGINS"); customOrigins != "" {
		corsOrigins = strings.Split(customOrigins, ",")
	}

	return cors.New(cors.Options{
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
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})
}