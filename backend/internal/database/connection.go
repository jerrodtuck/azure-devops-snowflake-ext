package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/snowflakedb/gosnowflake"
)

var DB *sql.DB

// InitializeDatabase establishes the Snowflake connection with retry logic
func InitializeDatabase() error {
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
		DB, err = sql.Open("snowflake", getConnectionString())
		if err != nil {
			return fmt.Errorf("failed to create database connection: %v", err)
		}

		// Test the connection
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err = DB.PingContext(ctx)
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
				if DB != nil {
					DB.Close()
					DB = nil
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

// getConnectionString builds the Snowflake connection string
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

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}