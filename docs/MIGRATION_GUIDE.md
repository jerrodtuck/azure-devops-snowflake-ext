# Migration Guide: From Hardcoded to Dynamic Queries

## Current State
Your `main.go` currently has hardcoded queries in two places:
1. `queries` map (lines 125-166) - for loading all data
2. `searchQueries` map (lines 276-298) - for search functionality

## Migration Steps

### Option 1: Complete Migration (Recommended)

1. **Update your main.go to use the dynamic system:**

```go
// Replace handleDropdownData with:
func handleDropdownData(w http.ResponseWriter, r *http.Request) {
    // Use the new dynamic handler
    handleSearch(w, r)
}

// Replace handleSearchData with:
func handleSearchData(w http.ResponseWriter, r *http.Request) {
    // Use the new dynamic handler
    handleSearch(w, r)
}

// Remove these variables:
// - var queries = map[string]string{...}
// - var searchQueries = map[string]string{...}
```

2. **Load configuration at startup:**

```go
func main() {
    // ... existing code ...
    
    // Load dynamic configuration
    if err := LoadConfig(); err != nil {
        log.Printf("Warning: Failed to load config: %v", err)
    }
    
    // ... rest of main ...
}
```

3. **Update routes to include new endpoints:**

```go
// Add these routes
router.HandleFunc("/api/config", handleGetConfig).Methods("GET")
router.HandleFunc("/api/search/{type}", handleSearch).Methods("GET")
```

### Option 2: Gradual Migration

Keep both systems running side-by-side:

```go
// In handleDropdownData, check if dynamic config exists first
func handleDropdownData(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    dataType := vars["type"]
    
    // Try dynamic config first
    if config, err := GetDataTypeConfig(dataType); err == nil {
        handleSearch(w, r)
        return
    }
    
    // Fall back to hardcoded queries
    // ... existing code ...
}
```

### Option 3: Keep as Override

Use hardcoded queries as defaults, but allow config.json to override:

```go
// In config.go, modify getDefaultConfig():
func getDefaultConfig() *Config {
    // Convert existing hardcoded queries to config format
    return &Config{
        DataTypes: []DataTypeConfig{
            {
                ID:   "cc",
                Name: "Cost Centers",
                Query: queries["cc"], // Use existing query
                SearchFields: []string{"COST_CENTER_NUMBER", "COST_CENTER_NAME"},
                Icon: "💰",
                Enabled: true,
            },
            // ... other types ...
        },
    }
}
```

## Benefits of Migration

1. **Customer Flexibility**: Each customer can have different data types
2. **No Code Changes**: Add new data types without recompiling
3. **Easier Testing**: Switch between configurations easily
4. **Better Maintenance**: All queries in one place (config.json)

## Example config.json

Create this file to replace your hardcoded queries. Copy from config.example.json:

```bash
cp config.example.json config.json
```

Then edit config.json to match your database schema:

```json
{
  "dataTypes": [
    {
      "id": "cc",
      "name": "Cost Centers",
      "description": "Company cost centers",
      "query": "SELECT id as value, id || ' - ' || name as label FROM your_database.your_schema.cost_centers WHERE (? = '' OR UPPER(name) LIKE UPPER('%' || ? || '%')) ORDER BY id LIMIT 100",
      "searchFields": ["name"],
      "icon": "💰",
      "enabled": true
    }
  ],
  "defaultDataType": "cc",
  "searchSettings": {
    "minSearchLength": 2,
    "debounceMs": 300,
    "maxResults": 100
  }
}
```

Note: config.json is gitignored to protect your sensitive database information.

## Testing the Migration

1. Create config.json with your queries
2. Run: `CONFIG_FILE=config.json go run *.go`
3. Test endpoints:
   - `GET /api/config` - Should return your data types
   - `GET /api/search/cc?q=1000` - Should search using dynamic query
4. Frontend will automatically adapt to available types