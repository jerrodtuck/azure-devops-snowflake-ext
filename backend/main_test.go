package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"snowflake-dropdown-api/internal/models"

	"github.com/gorilla/mux"
)

// Mock data for testing without Snowflake connection
var mockCCData = []models.DropdownItem{
	{Value: "CC001", Label: "CC001 - Marketing SHOWBOAT"},
	{Value: "CC002", Label: "CC002 - Engineering"},
	{Value: "CC003", Label: "CC003 - Sales SHOWBOAT"},
	{Value: "CC004", Label: "CC004 - Finance"},
	{Value: "CC005", Label: "CC005 - HR Resources"},
}

var mockWBSData = []models.DropdownItem{
	{Value: "WBS001", Label: "WBS001 - Project BONSAI Alpha"},
	{Value: "WBS002", Label: "WBS002 - Project Beta"},
	{Value: "WBS003", Label: "WBS003 - BONSAI Development"},
	{Value: "WBS004", Label: "WBS004 - Infrastructure"},
	{Value: "WBS005", Label: "WBS005 - Research BONSAI"},
}

// Test mode handler
func handleMockSearch(w http.ResponseWriter, r *http.Request) {
	// Validate request method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	dataType := vars["type"]

	// Validate dataType parameter
	if dataType == "" {
		http.Error(w, "Missing type parameter", http.StatusBadRequest)
		return
	}

	// Sanitize and validate search term
	searchTerm := strings.TrimSpace(r.URL.Query().Get("q"))
	if len(searchTerm) > 100 { // Reasonable limit
		http.Error(w, "Search term too long", http.StatusBadRequest)
		return
	}

	var results []models.DropdownItem
	var sourceData []models.DropdownItem

	// Select source data based on type
	switch strings.ToLower(dataType) {
	case "cc":
		sourceData = mockCCData
	case "wbs":
		sourceData = mockWBSData
	default:
		http.Error(w, fmt.Sprintf("Unknown type: %s", dataType), http.StatusBadRequest)
		return
	}

	// Filter by search term
	if searchTerm == "" {
		results = sourceData
	} else {
		for _, item := range sourceData {
			if containsIgnoreCase(item.Label, searchTerm) || containsIgnoreCase(item.Value, searchTerm) {
				results = append(results, item)
			}
		}
	}

	// Return response
	response := models.DropdownResponse{
		Data: results,
		Metadata: models.Metadata{
			ExportedAt: time.Now().UTC(),
			RowCount:   len(results),
			Source:     dataType + "_mock",
			Cached:     false,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Helper function for case-insensitive contains
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToUpper(s), strings.ToUpper(substr))
}

// Add test mode to main function
func setupTestMode(router *mux.Router) {
	// Override search endpoint with mock data
	router.HandleFunc("/api/search/{type}", handleMockSearch).Methods("GET")

	// Add test data endpoint
	router.HandleFunc("/api/test/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		response := map[string]interface{}{
			"cc_data":  mockCCData,
			"wbs_data": mockWBSData,
			"message":  "Test mode active - using mock data",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode test data response", http.StatusInternalServerError)
			return
		}
	}).Methods("GET")
}
