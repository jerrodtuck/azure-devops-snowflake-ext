package database

import (
	"strings"
	"snowflake-dropdown-api/internal/config"
)

// BuildSearchQuery builds a parameterized query for searching
func BuildSearchQuery(dtConfig *config.DataTypeConfig, searchTerm string) (string, []interface{}) {
	params := []interface{}{searchTerm}

	// Add parameters for each search field
	for range dtConfig.SearchFields {
		params = append(params, "%"+strings.ToUpper(searchTerm)+"%")
	}

	// Add the limit parameter from global config
	maxResults := 100 // default
	if appConfig := config.AppConfig; appConfig != nil {
		maxResults = appConfig.SearchSettings.MaxResults
	}
	params = append(params, maxResults)

	return dtConfig.Query, params
}