package handlers

import (
	"encoding/json"
	"net/http"
	"snowflake-dropdown-api/internal/config"
	"snowflake-dropdown-api/internal/models"
)

// HandleGetConfig returns the configuration for the frontend
func HandleGetConfig(w http.ResponseWriter, r *http.Request) {
	enabledTypes := config.GetEnabledDataTypes()

	var dataTypes []models.DataTypeInfo
	for _, dt := range enabledTypes {
		dataTypes = append(dataTypes, models.DataTypeInfo{
			ID:          dt.ID,
			Name:        dt.Name,
			Description: dt.Description,
			Icon:        dt.Icon,
		})
	}

	response := models.ConfigResponse{
		DataTypes:   dataTypes,
		DefaultType: config.AppConfig.DefaultDataType,
		SearchSettings: struct {
			MinSearchLength int `json:"minSearchLength"`
			DebounceMs      int `json:"debounceMs"`
		}{
			MinSearchLength: config.AppConfig.SearchSettings.MinSearchLength,
			DebounceMs:      config.AppConfig.SearchSettings.DebounceMs,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleGetDataTypes returns available data types
func HandleGetDataTypes(w http.ResponseWriter, r *http.Request) {
	types := config.GetEnabledDataTypes()

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