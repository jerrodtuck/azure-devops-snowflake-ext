# Snowflake Dropdown API - Go Backend

A fast, efficient Go server that provides dropdown data from Snowflake for Azure DevOps.

## Features

- ✅ **Fast Performance** - Go's concurrency and efficiency
- ✅ **Built-in Caching** - 1-hour cache to reduce Snowflake queries
- ✅ **CORS Support** - Works with Azure DevOps
- ✅ **Multiple Datasets** - Support for cost centers, departments, products
- ✅ **Health Checks** - For monitoring
- ✅ **Docker Ready** - Easy deployment

## Quick Start

### 1. Local Development

```bash
# Install dependencies
go mod download

# Set environment variables
cp .env.example .env
# Edit .env with your Snowflake credentials

# Copy and configure data types
cp config.example.json config.json
# Edit config.json with your database schema and queries

# Run the server
go run main.go
```

### 2. Test the API

```bash
# Health check
curl http://localhost:8080/api/health

# Get cost centers
curl http://localhost:8080/api/dropdown/cost_centers

# Get available data types
curl http://localhost:8080/api/types
```

## API Endpoints

| Endpoint | Description | Response |
|----------|-------------|----------|
| `GET /api/health` | Health check | `{"status": "healthy"}` |
| `GET /api/dropdown/{type}` | Get dropdown data | Dropdown items with metadata |
| `GET /api/types` | List available types | `{"available_types": [...]}` |

## Response Format

```json
{
  "data": [
    {
      "value": "CC001",
      "label": "CC001 - Marketing Department"
    },
    {
      "value": "CC002",
      "label": "CC002 - Engineering"
    }
  ],
  "metadata": {
    "exported_at": "2024-01-15T10:30:00Z",
    "row_count": 150,
    "source": "cost_centers",
    "cached": false
  }
}
```

## Deployment Options

### Option 1: Docker

```bash
# Build
docker build -t snowflake-dropdown-api .

# Run
docker run -p 8080:8080 --env-file .env snowflake-dropdown-api
```

### Option 2: Azure Container Instances

```bash
# Build and push to Azure Container Registry
az acr build --registry myregistry --image snowflake-api .

# Deploy
az container create \
  --resource-group mygroup \
  --name snowflake-api \
  --image myregistry.azurecr.io/snowflake-api \
  --dns-name-label snowflake-api \
  --ports 8080 \
  --environment-variables-file .env
```

### Option 3: Azure App Service

```bash
# Deploy to App Service
az webapp up --name snowflake-api --runtime "GO:1.21"
```

## Configuration

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `SNOWFLAKE_ACCOUNT` | Snowflake account | `your-account` |
| `SNOWFLAKE_USER` | Username | `SERVICE_ACCOUNT` |
| `SNOWFLAKE_PASSWORD` | Password | `****` |
| `SNOWFLAKE_DATABASE` | Database name | `FINANCE_AND_ACCOUNTING` |
| `SNOWFLAKE_SCHEMA` | Schema name | `GOLDEN` |
| `SNOWFLAKE_WAREHOUSE` | Warehouse | `DATABRICKS` |
| `SNOWFLAKE_ROLE` | Role | `your-role` |
| `PORT` | Server port | `8080` |

### Adding New Data Types

Edit `main.go` and add to the `queries` map:

```go
var queries = map[string]string{
    "your_new_type": `
        SELECT 
            ID as value,
            NAME as label
        FROM YOUR_TABLE
        WHERE ACTIVE = TRUE
    `,
}
```

## Security Considerations

1. **Use HTTPS** in production
2. **Implement authentication** (JWT, API keys)
3. **Rate limiting** to prevent abuse
4. **Network restrictions** - whitelist Azure DevOps IPs
5. **Use service accounts** with minimal permissions

## Performance

- Queries are cached for 1 hour
- Concurrent requests are handled efficiently
- Connection pooling for Snowflake
- Typical response time: <100ms (cached), <2s (fresh query)

## Monitoring

The `/api/health` endpoint can be used for:
- Azure Monitor
- Application Insights
- Load balancer health checks

## Extension Configuration

In your Azure DevOps extension, configure:

```javascript
// For production
DataUrl: "https://your-api.azurewebsites.net/api/dropdown/cost_centers"

// For local testing
DataUrl: "http://localhost:8080/api/dropdown/cost_centers"
```