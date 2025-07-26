package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

var AppConfig *Config

// LoadConfig loads configuration from file or environment
func LoadConfig() error {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.json"
	}

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Printf("Config file %s not found, using default configuration", configFile)
		AppConfig = getDefaultConfig()
		return nil
	}

	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("error parsing config file: %v", err)
	}

	AppConfig = &config
	log.Printf("Loaded configuration with %d data types", len(config.DataTypes))
	return nil
}

// LoadEnvFile loads environment variables from .env file
func LoadEnvFile() error {
	file, err := os.Open(".env")
	if err != nil {
		// .env file is optional
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first = sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Only set if not already set (allows overriding with actual env vars)
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
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

// getDefaultConfig returns the default configuration
func getDefaultConfig() *Config {
	return &Config{
		DataTypes: []DataTypeConfig{
			{
				ID:           "cc",
				Name:         "Cost Centers",
				Description:  "Company cost centers",
				Query:        "SELECT COST_CENTER_NUMBER as value, COST_CENTER_NUMBER || ' - ' || COST_CENTER_NAME as label FROM FINANCE_AND_ACCOUNTING.golden.cost_center WHERE (? = '' OR COST_CENTER_NUMBER LIKE ? OR COST_CENTER_NAME LIKE ?) ORDER BY COST_CENTER_NUMBER LIMIT 100",
				SearchFields: []string{"COST_CENTER_NUMBER", "COST_CENTER_NAME"},
				Icon:         "💰",
				Enabled:      true,
			},
			{
				ID:           "wbs",
				Name:         "WBS Elements",
				Description:  "Work Breakdown Structure elements",
				Query:        "SELECT WBS_NUMBER as value, WBS_NUMBER || ' - ' || WBS_DESCRIPTION as label FROM FINANCE_AND_ACCOUNTING.golden.wbs WHERE (? = '' OR WBS_NUMBER LIKE ? OR WBS_DESCRIPTION LIKE ?) ORDER BY WBS_NUMBER LIMIT 100",
				SearchFields: []string{"WBS_NUMBER", "WBS_DESCRIPTION"},
				Icon:         "📊",
				Enabled:      true,
			},
		},
		DefaultDataType: "cc",
		CacheSettings: struct {
			Enabled    bool `json:"enabled"`
			TTLMinutes int  `json:"ttlMinutes"`
		}{
			Enabled:    true,
			TTLMinutes: 60,
		},
		SearchSettings: struct {
			MinSearchLength int `json:"minSearchLength"`
			DebounceMs      int `json:"debounceMs"`
			MaxResults      int `json:"maxResults"`
		}{
			MinSearchLength: 2,
			DebounceMs:      300,
			MaxResults:      100,
		},
	}
}

// GetDataTypeConfig returns configuration for a specific data type
func GetDataTypeConfig(dataType string) (*DataTypeConfig, error) {
	if AppConfig == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}

	for _, dt := range AppConfig.DataTypes {
		if dt.ID == dataType && dt.Enabled {
			return &dt, nil
		}
	}

	return nil, fmt.Errorf("data type '%s' not found or disabled", dataType)
}

// GetEnabledDataTypes returns all enabled data types
func GetEnabledDataTypes() []DataTypeConfig {
	if AppConfig == nil {
		return []DataTypeConfig{}
	}

	var enabled []DataTypeConfig
	for _, dt := range AppConfig.DataTypes {
		if dt.Enabled {
			enabled = append(enabled, dt)
		}
	}
	return enabled
}