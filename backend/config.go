package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// DataTypeConfig represents a configurable data type
type DataTypeConfig struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Query        string   `json:"query"`
	SearchFields []string `json:"searchFields"`
	Icon         string   `json:"icon"`
	Enabled      bool     `json:"enabled"`
}

// Config represents the application configuration
type Config struct {
	DataTypes      []DataTypeConfig `json:"dataTypes"`
	DefaultDataType string          `json:"defaultDataType"`
	CacheSettings  struct {
		Enabled    bool `json:"enabled"`
		TTLMinutes int  `json:"ttlMinutes"`
	} `json:"cacheSettings"`
	SearchSettings struct {
		MinSearchLength int `json:"minSearchLength"`
		DebounceMs      int `json:"debounceMs"`
		MaxResults      int `json:"maxResults"`
	} `json:"searchSettings"`
}

var appConfig *Config

// LoadConfig loads configuration from file or environment
func LoadConfig() error {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.json"
	}

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Printf("Config file %s not found, using default configuration", configFile)
		appConfig = getDefaultConfig()
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

	appConfig = &config
	log.Printf("Loaded configuration with %d data types", len(config.DataTypes))
	return nil
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
	if appConfig == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}

	for _, dt := range appConfig.DataTypes {
		if dt.ID == dataType && dt.Enabled {
			return &dt, nil
		}
	}

	return nil, fmt.Errorf("data type '%s' not found or disabled", dataType)
}

// GetEnabledDataTypes returns all enabled data types
func GetEnabledDataTypes() []DataTypeConfig {
	if appConfig == nil {
		return []DataTypeConfig{}
	}

	var enabled []DataTypeConfig
	for _, dt := range appConfig.DataTypes {
		if dt.Enabled {
			enabled = append(enabled, dt)
		}
	}
	return enabled
}

// BuildSearchQuery builds a parameterized query for searching
func BuildSearchQuery(config *DataTypeConfig, searchTerm string) (string, []interface{}) {
	params := []interface{}{searchTerm}
	
	// Add parameters for each search field
	for range config.SearchFields {
		params = append(params, "%"+strings.ToUpper(searchTerm)+"%")
	}
	
	return config.Query, params
}