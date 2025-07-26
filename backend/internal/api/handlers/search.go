package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"snowflake-dropdown-api/internal/config"
	"snowflake-dropdown-api/internal/database"
	"snowflake-dropdown-api/internal/models"

	"github.com/gorilla/mux"
)

// HandleSearch performs a search using the dynamic configuration
func HandleSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dataType := vars["type"]

	// Validate dataType parameter
	if dataType == "" {
		http.Error(w, "Data type is required", http.StatusBadRequest)
		return
	}

	searchTerm := r.URL.Query().Get("q")

	// Handle test mode
	if os.Getenv("TEST_MODE") == "true" {
		handleMockSearch(w, dataType, searchTerm)
		return
	}

	// Get configuration for this data type
	dtConfig, err := config.GetDataTypeConfig(dataType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Build and execute query
	query, params := database.BuildSearchQuery(dtConfig, searchTerm)

	rows, err := database.DB.Query(query, params...)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.DropdownItem
	for rows.Next() {
		var item models.DropdownItem
		if err := rows.Scan(&item.Value, &item.Label); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		items = append(items, item)
	}

	response := models.DropdownResponse{
		Data: items,
		Metadata: models.Metadata{
			ExportedAt: time.Now().UTC(),
			RowCount:   len(items),
			Source:     dataType,
			Cached:     false,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleDynamicSearch handles custom queries with security validation
func HandleDynamicSearch(w http.ResponseWriter, r *http.Request) {
	var request models.DynamicSearchRequest

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

	rows, err := database.DB.Query(request.Query, params...)
	if err != nil {
		log.Printf("Dynamic query error: %v", err)
		http.Error(w, "Query execution failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.DropdownItem
	for rows.Next() {
		var item models.DropdownItem
		if err := rows.Scan(&item.Value, &item.Label); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.DropdownResponse{
		Data: items,
		Metadata: models.Metadata{
			ExportedAt: time.Now().UTC(),
			RowCount:   len(items),
			Source:     "dynamic",
			Cached:     false,
		},
	})
}

// handleMockSearch returns mock data for testing
func handleMockSearch(w http.ResponseWriter, dataType, searchTerm string) {
	var items []models.DropdownItem

	switch dataType {
	case "cc", "cost_centers":
		items = []models.DropdownItem{
			{Value: "1000", Label: "1000 - IT Department"},
			{Value: "2000", Label: "2000 - Finance Department"},
			{Value: "3000", Label: "3000 - Marketing Department"},
			{Value: "4000", Label: "4000 - Operations"},
			{Value: "5000", Label: "5000 - Human Resources"},
		}
	case "wbs":
		items = []models.DropdownItem{
			{Value: "WBS001", Label: "WBS001 - Project Alpha"},
			{Value: "WBS002", Label: "WBS002 - Project Beta"},
			{Value: "WBS003", Label: "WBS003 - Project Gamma"},
			{Value: "WBS004", Label: "WBS004 - Project Delta"},
			{Value: "WBS005", Label: "WBS005 - Project Epsilon"},
		}
	default:
		http.Error(w, fmt.Sprintf("Unknown data type: %s", dataType), http.StatusBadRequest)
		return
	}

	// Filter by search term if provided
	if searchTerm != "" {
		searchLower := strings.ToLower(searchTerm)
		var filtered []models.DropdownItem

		for _, item := range items {
			if strings.Contains(strings.ToLower(item.Value), searchLower) ||
				strings.Contains(strings.ToLower(item.Label), searchLower) {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}

	response := models.DropdownResponse{
		Data: items,
		Metadata: models.Metadata{
			ExportedAt: time.Now().UTC(),
			RowCount:   len(items),
			Source:     dataType + " (mock)",
			Cached:     false,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
