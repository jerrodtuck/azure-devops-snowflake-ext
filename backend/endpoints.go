package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// This file expects the following globals from main.go:
// - db (*sql.DB): Database connection
// - appConfig (*Config): Application configuration from config.go
// - DropdownItem, DropdownResponse, and Metadata types

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

// handleGetConfig returns the configuration for the frontend
func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	enabledTypes := GetEnabledDataTypes()
	
	var dataTypes []DataTypeInfo
	for _, dt := range enabledTypes {
		dataTypes = append(dataTypes, DataTypeInfo{
			ID:          dt.ID,
			Name:        dt.Name,
			Description: dt.Description,
			Icon:        dt.Icon,
		})
	}

	response := ConfigResponse{
		DataTypes:   dataTypes,
		DefaultType: appConfig.DefaultDataType,
		SearchSettings: struct {
			MinSearchLength int `json:"minSearchLength"`
			DebounceMs      int `json:"debounceMs"`
		}{
			MinSearchLength: appConfig.SearchSettings.MinSearchLength,
			DebounceMs:      appConfig.SearchSettings.DebounceMs,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSearch performs a search using the dynamic configuration
func handleSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dataType := vars["type"]
	searchTerm := r.URL.Query().Get("q")

	// Get configuration for this data type
	config, err := GetDataTypeConfig(dataType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Build and execute query
	query, params := BuildSearchQuery(config, searchTerm)
	
	rows, err := db.Query(query, params...)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []DropdownItem
	for rows.Next() {
		var item DropdownItem
		if err := rows.Scan(&item.Value, &item.Label); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		items = append(items, item)
	}

	response := DropdownResponse{
		Data: items,
		Metadata: Metadata{
			ExportedAt: time.Now().UTC(),
			RowCount:   len(items),
			Source:     dataType,
			Cached:     false,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleGetDataTypes returns available data types
func handleGetDataTypes(w http.ResponseWriter, r *http.Request) {
	types := GetEnabledDataTypes()
	
	var response []map[string]string
	for _, t := range types {
		response = append(response, map[string]string{
			"id":   t.ID,
			"name": t.Name,
			"icon": t.Icon,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Dynamic search handler that supports custom queries
func handleDynamicSearch(w http.ResponseWriter, r *http.Request) {
	// Parse request body for custom query
	var request struct {
		Query      string   `json:"query"`
		Parameters []string `json:"parameters"`
		DataType   string   `json:"dataType"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Security check - only allow SELECT queries
	if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(request.Query)), "SELECT") {
		http.Error(w, "Only SELECT queries are allowed", http.StatusForbidden)
		return
	}

	// Execute query with parameters
	params := make([]interface{}, len(request.Parameters))
	for i, p := range request.Parameters {
		params[i] = p
	}

	rows, err := db.Query(request.Query, params...)
	if err != nil {
		log.Printf("Dynamic query error: %v", err)
		http.Error(w, "Query execution failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []DropdownItem
	for rows.Next() {
		var item DropdownItem
		if err := rows.Scan(&item.Value, &item.Label); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DropdownResponse{
		Data: items,
		Metadata: Metadata{
			ExportedAt: time.Now().UTC(),
			RowCount:   len(items),
			Source:     "dynamic",
			Cached:     false,
		},
	})
}