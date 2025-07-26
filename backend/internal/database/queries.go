package database

import (
	"strings"
	"snowflake-dropdown-api/internal/config"
)

// BuildSearchQuery builds a parameterized query for searching
func BuildSearchQuery(config *config.DataTypeConfig, searchTerm string) (string, []interface{}) {
	params := []interface{}{searchTerm}

	// Add parameters for each search field
	for range config.SearchFields {
		params = append(params, "%"+strings.ToUpper(searchTerm)+"%")
	}

	return config.Query, params
}