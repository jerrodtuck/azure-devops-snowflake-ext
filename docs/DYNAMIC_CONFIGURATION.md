# Dynamic Configuration Guide

This guide explains how to make the Azure DevOps Snowflake Extension dynamic and customizable for different customers.

## Overview

The extension now supports dynamic data types instead of hardcoded "Cost Center (CC)" and "WBS Element" options. You can configure:
- Custom data types (e.g., Projects, Departments, Employees, Vendors)
- Custom SQL queries for each data type
- Search behavior and UI settings
- Icons and descriptions

## Configuration File

Create a `config.json` file in your backend directory:

```json
{
  "dataTypes": [
    {
      "id": "projects",
      "name": "Projects",
      "description": "Active projects",
      "query": "SELECT PROJECT_ID as value, PROJECT_ID || ' - ' || PROJECT_NAME as label FROM PROJECTS WHERE STATUS = 'ACTIVE' AND (? = '' OR PROJECT_ID LIKE ? OR PROJECT_NAME LIKE ?) ORDER BY PROJECT_ID LIMIT 100",
      "searchFields": ["PROJECT_ID", "PROJECT_NAME"],
      "icon": "📁",
      "enabled": true
    },
    {
      "id": "vendors",
      "name": "Vendors",
      "description": "Approved vendors",
      "query": "SELECT VENDOR_ID as value, VENDOR_NAME as label FROM VENDORS WHERE APPROVED = true AND (? = '' OR VENDOR_ID LIKE ? OR VENDOR_NAME LIKE ?) ORDER BY VENDOR_NAME LIMIT 100",
      "searchFields": ["VENDOR_ID", "VENDOR_NAME"],
      "icon": "🏢",
      "enabled": true
    }
  ],
  "defaultDataType": "projects",
  "searchSettings": {
    "minSearchLength": 2,
    "debounceMs": 300,
    "maxResults": 100
  }
}
```

## Query Requirements

Each query must:
1. Return two columns: `value` and `label`
2. Accept search parameters for filtering
3. Include proper WHERE clauses for search
4. Have ORDER BY for consistent results
5. Include LIMIT to prevent large result sets

### Query Template
```sql
SELECT 
  YOUR_ID_COLUMN as value,
  YOUR_ID_COLUMN || ' - ' || YOUR_DESCRIPTION_COLUMN as label
FROM YOUR_TABLE
WHERE (? = '' OR YOUR_ID_COLUMN LIKE ? OR YOUR_DESCRIPTION_COLUMN LIKE ?)
ORDER BY YOUR_ID_COLUMN
LIMIT 100
```

## Implementation Steps

### 1. Backend Configuration

1. Copy `config.json.example` to `config.json`
2. Customize data types for your customer
3. Update the Go backend to use dynamic configuration:

```go
// In main.go, add after database initialization:
if err := LoadConfig(); err != nil {
    log.Printf("Warning: Failed to load config: %v", err)
}

// Add new routes:
api.HandleFunc("/config", handleGetConfig).Methods("GET", "OPTIONS")
api.HandleFunc("/search/{type}", handleSearch).Methods("GET", "OPTIONS")
```

### 2. Frontend Updates

The frontend automatically fetches available data types from the `/api/config` endpoint. No code changes needed!

### 3. Azure DevOps Configuration

In the extension settings:
- **API URL**: Your backend URL (e.g., `https://api.customer.com/api`)
- **Data Type**: Will be dynamically populated based on config

## Customer-Specific Configurations

### Example: Financial Services
```json
{
  "dataTypes": [
    {
      "id": "cost_center",
      "name": "Cost Centers",
      "query": "SELECT CC_CODE as value, CC_CODE || ' - ' || CC_NAME as label FROM COST_CENTERS WHERE ACTIVE = 'Y'",
      "icon": "💰"
    },
    {
      "id": "gl_account",
      "name": "GL Accounts",
      "query": "SELECT ACCOUNT_NO as value, ACCOUNT_NO || ' - ' || ACCOUNT_DESC as label FROM GL_ACCOUNTS WHERE ACTIVE = 1",
      "icon": "📊"
    }
  ]
}
```

### Example: Manufacturing
```json
{
  "dataTypes": [
    {
      "id": "plant",
      "name": "Plants",
      "query": "SELECT PLANT_CODE as value, PLANT_CODE || ' - ' || PLANT_NAME as label FROM PLANTS WHERE STATUS = 'OPERATIONAL'",
      "icon": "🏭"
    },
    {
      "id": "material",
      "name": "Materials",
      "query": "SELECT MATERIAL_NO as value, MATERIAL_NO || ' - ' || DESCRIPTION as label FROM MATERIALS WHERE ACTIVE = true",
      "icon": "📦"
    }
  ]
}
```

### Example: Healthcare
```json
{
  "dataTypes": [
    {
      "id": "department",
      "name": "Departments",
      "query": "SELECT DEPT_ID as value, DEPT_ID || ' - ' || DEPT_NAME as label FROM DEPARTMENTS WHERE ACTIVE = 1",
      "icon": "🏥"
    },
    {
      "id": "provider",
      "name": "Providers",
      "query": "SELECT PROVIDER_ID as value, PROVIDER_ID || ' - ' || PROVIDER_NAME as label FROM PROVIDERS WHERE STATUS = 'ACTIVE'",
      "icon": "👨‍⚕️"
    }
  ]
}
```

## Security Considerations

1. **Query Validation**: Only SELECT queries are allowed
2. **Parameterized Queries**: Prevents SQL injection
3. **Result Limits**: Enforced to prevent large data transfers
4. **CORS**: Configured for Azure DevOps domains only
5. **Environment-based Config**: Use different configs for dev/prod

## Deployment

1. **Development**:
   ```bash
   CONFIG_FILE=config.dev.json go run main.go
   ```

2. **Production**:
   ```bash
   CONFIG_FILE=config.prod.json ./snowflake-backend
   ```

3. **Docker**:
   ```dockerfile
   COPY config.json /app/config.json
   ENV CONFIG_FILE=/app/config.json
   ```

## Troubleshooting

1. **Data types not showing**: Check `/api/config` endpoint
2. **Search not working**: Verify query syntax and parameters
3. **Wrong results**: Check ORDER BY and search conditions
4. **Performance issues**: Add database indexes on search columns

## Future Enhancements

1. **Admin UI**: Web interface for managing configurations
2. **Query Builder**: Visual tool for creating queries
3. **Multi-tenancy**: Different configs per Azure DevOps organization
4. **Caching**: Redis integration for better performance
5. **Audit Logging**: Track configuration changes