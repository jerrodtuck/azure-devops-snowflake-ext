package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Mock data for testing without Snowflake connection
var mockCCData = []DropdownItem{
	{Value: "CC001", Label: "CC001 - Marketing SHOWBOAT"},
	{Value: "CC002", Label: "CC002 - Engineering"},
	{Value: "CC003", Label: "CC003 - Sales SHOWBOAT"},
	{Value: "CC004", Label: "CC004 - Finance"},
	{Value: "CC005", Label: "CC005 - HR Resources"},
}

var mockWBSData = []DropdownItem{
	{Value: "WBS001", Label: "WBS001 - Project BONSAI Alpha"},
	{Value: "WBS002", Label: "WBS002 - Project Beta"},
	{Value: "WBS003", Label: "WBS003 - BONSAI Development"},
	{Value: "WBS004", Label: "WBS004 - Infrastructure"},
	{Value: "WBS005", Label: "WBS005 - Research BONSAI"},
}

// Test mode handler
func handleMockSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dataType := vars["type"]
	searchTerm := r.URL.Query().Get("q")

	var results []DropdownItem
	var sourceData []DropdownItem

	// Select source data based on type
	switch dataType {
	case "cc":
		sourceData = mockCCData
	case "wbs":
		sourceData = mockWBSData
	default:
		http.Error(w, "Unknown type", http.StatusBadRequest)
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
	response := DropdownResponse{
		Data: results,
		Metadata: Metadata{
			ExportedAt: time.Now().UTC(),
			RowCount:   len(results),
			Source:     dataType + "_mock",
			Cached:     false,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
		json.NewEncoder(w).Encode(map[string]interface{}{
			"cc_data":  mockCCData,
			"wbs_data": mockWBSData,
			"message":  "Test mode active - using mock data",
		})
	}).Methods("GET")
}
