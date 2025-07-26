package models

import "time"

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

// DataTypeInfo represents information about a data type for the frontend
type DataTypeInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// ConfigResponse represents the configuration response for the frontend
type ConfigResponse struct {
	DataTypes      []DataTypeInfo `json:"dataTypes"`
	DefaultType    string         `json:"defaultType"`
	SearchSettings struct {
		MinSearchLength int `json:"minSearchLength"`
		DebounceMs      int `json:"debounceMs"`
	} `json:"searchSettings"`
}

// DynamicSearchRequest represents a request for dynamic search
type DynamicSearchRequest struct {
	Query      string   `json:"query"`
	Parameters []string `json:"parameters"`
	DataType   string   `json:"dataType"`
}