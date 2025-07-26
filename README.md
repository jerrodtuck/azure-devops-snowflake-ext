# Azure DevOps Snowflake Extension

A modern Azure DevOps extension that provides searchable dropdown controls populated with data from Snowflake.

## Project Structure

- **extension/** - React-based Azure DevOps extension
- **backend/go-backend/** - Go API server that connects to Snowflake
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
   `ash
   cd backend/go-backend
   `

2. Copy .env.example to .env and configure:
   `
   SNOWFLAKE_ACCOUNT=your-account
   SNOWFLAKE_USER=your-user
   SNOWFLAKE_PASSWORD=your-password
   SNOWFLAKE_DATABASE=your-db
   SNOWFLAKE_SCHEMA=your-schema
   SNOWFLAKE_WAREHOUSE=your-warehouse
   SNOWFLAKE_ROLE=your-role
   `

3. Run the server:
   `ash
   go run main.go loadenv.go
   `

### Extension Setup

1. Navigate to the extension directory:
   `ash
   cd extension
   `

2. Install dependencies:
   `ash
   npm install
   `

3. Build the extension:
   `ash
   npm run build
   `

4. Package for Azure DevOps:
   `ash
   npm run package
   `

## Development

### Backend Development
`ash
cd backend/go-backend
go run main.go loadenv.go
`

### Extension Development
`ash
cd extension
npm start
`

## Deployment

See [docs/](docs/) for detailed deployment guides.

## License

[Your License]
