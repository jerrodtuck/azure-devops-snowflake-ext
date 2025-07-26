# Integration Guide for Dynamic Configuration

This guide shows how to integrate the new dynamic configuration files into your existing main.go.

## Files Added

1. **config.go** - Configuration management system
2. **endpoints.go** - New API endpoints for dynamic data types
3. **config.json.example** - Example configuration file

## Integration Steps

### 1. Add Configuration Loading to main.go

In your `main()` function, after database initialization:

```go
// Load configuration
if err := LoadConfig(); err != nil {
    log.Printf("Warning: Failed to load config: %v", err)
    // Continue with defaults
}
```

### 2. Add New Routes

Replace or update your existing routes with:

```go
// API routes
api := router.PathPrefix("/api").Subrouter()

// Health check (existing)
api.HandleFunc("/health", handleHealth).Methods("GET", "OPTIONS")

// NEW: Configuration endpoint
api.HandleFunc("/config", handleGetConfig).Methods("GET", "OPTIONS")

// NEW: Dynamic search (replaces old dropdown endpoints)
api.HandleFunc("/search/{type}", handleSearch).Methods("GET", "OPTIONS")

// NEW: List data types
api.HandleFunc("/types", handleGetDataTypes).Methods("GET", "OPTIONS")

// Keep your existing endpoints for backward compatibility
api.HandleFunc("/dropdown/{type}", handleDropdown).Methods("GET", "OPTIONS")
```

### 3. Update handleDropdown (Optional)

To maintain backward compatibility while using new system:

```go
func handleDropdown(w http.ResponseWriter, r *http.Request) {
    // Use the new search handler
    handleSearch(w, r)
}
```

### 4. Create config.json

Copy `config.json.example` to `config.json` and customize:

```json
{
  "dataTypes": [
    {
      "id": "cc",
      "name": "Cost Centers",
      "description": "Company cost centers",
      "query": "SELECT COST_CENTER_NUMBER as value, COST_CENTER_NUMBER || ' - ' || COST_CENTER_NAME as label FROM FINANCE_AND_ACCOUNTING.golden.cost_center WHERE (? = '' OR COST_CENTER_NUMBER LIKE ? OR COST_CENTER_NAME LIKE ?) ORDER BY COST_CENTER_NUMBER LIMIT 100",
      "searchFields": ["COST_CENTER_NUMBER", "COST_CENTER_NAME"],
      "icon": "💰",
      "enabled": true
    }
  ],
  "defaultDataType": "cc"
}
```

## Testing

1. Start the backend with config:
   ```bash
   CONFIG_FILE=config.json go run *.go
   ```

2. Test the new endpoints:
   ```bash
   # Get configuration
   curl http://localhost:8080/api/config
   
   # Search with dynamic type
   curl http://localhost:8080/api/search/cc?q=1000
   
   # List available types
   curl http://localhost:8080/api/types
   ```

## Notes

- The old hardcoded queries in main.go can remain for backward compatibility
- The new system will use config.json if available, otherwise falls back to defaults
- No frontend changes needed - it will automatically use the new endpoints