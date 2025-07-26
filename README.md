# Azure DevOps Snowflake Extension

A modern Azure DevOps extension that provides searchable dropdown controls populated with data from Snowflake.

## Project Structure

- **extension/** - React-based Azure DevOps extension
- **backend/** - Go API server that connects to Snowflake
- **docs/** - Documentation and guides

## Features

- 🔍 Live search with debouncing
- 📊 Data from Snowflake (Cost Centers, WBS Elements)
- ⚡ Fast Go backend with caching
- 🔒 Secure authentication (SSO support)
- ⌨️ Keyboard navigation
- 🎨 Modern React/TypeScript implementation

## Quick Start

### Prerequisites

- Node.js (v16+)
- Go (v1.19+)
- Azure DevOps organization
- Snowflake account

### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Copy .env.example to .env and configure:
   ```
   SNOWFLAKE_ACCOUNT=your-account
   SNOWFLAKE_USER=your-user
   SNOWFLAKE_PASSWORD=your-password
   SNOWFLAKE_DATABASE=your-db
   SNOWFLAKE_SCHEMA=your-schema
   SNOWFLAKE_WAREHOUSE=your-warehouse
   SNOWFLAKE_ROLE=your-role
   ```

3. Run the server:
   ```bash
   go run main.go loadenv.go
   ```

### Extension Setup

1. Navigate to the extension directory:
   ```bash
   cd extension
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Build the extension:
   ```bash
   npm run build
   ```

4. Package for Azure DevOps:
   ```bash
   npm run package
   ```

## Development

### Backend Development
```bash
cd backend
go run main.go loadenv.go config.go security.go endpoints.go
```

### Extension Development
```bash
cd extension
npm start
```

## Deployment

### Local Development
See [docs/](docs/) for detailed development guides.

### Azure Deployment Options

The backend can be deployed to Azure using several approaches:

- **Azure Container Instances (ACI)** - Simple container deployment
- **Azure App Service** - Managed platform with built-in scaling
- **Azure Kubernetes Service (AKS)** - Full container orchestration
- **Azure Container Apps** - Serverless container platform

### Docker Deployment

Docker configuration files are available in `backend/deployments/`:

- `Dockerfile` - Container image definition
- `docker-compose.yml` - Multi-container setup with dependencies

To deploy with Docker:
```bash
cd backend/deployments
docker-compose up -d
```
