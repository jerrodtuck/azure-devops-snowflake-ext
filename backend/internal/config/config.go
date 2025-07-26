package config

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
	DataTypes       []DataTypeConfig `json:"dataTypes"`
	DefaultDataType string           `json:"defaultDataType"`
	CacheSettings   struct {
		Enabled    bool `json:"enabled"`
		TTLMinutes int  `json:"ttlMinutes"`
	} `json:"cacheSettings"`
	SearchSettings struct {
		MinSearchLength int `json:"minSearchLength"`
		DebounceMs      int `json:"debounceMs"`
		MaxResults      int `json:"maxResults"`
	} `json:"searchSettings"`
}

// SecurityConfig holds security settings
type SecurityConfig struct {
	APIKeyEnabled bool
	APIKeys       []string
	JWTEnabled    bool
	JWTSecret     string
	IPWhitelist   []string
}