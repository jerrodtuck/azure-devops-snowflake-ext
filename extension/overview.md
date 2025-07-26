# Snowflake Dropdown Extension

A modern work item control that provides searchable dropdown fields populated with live data from Snowflake. You will need to implement the Go Backend or modify the project to use the Snowflake SQL API. For custom implementations you can reach out to me via the various contact methods on https://jerrodtuck.com

![Snowflake Dropdown Demo](https://raw.githubusercontent.com/jerrodtuck/azure-devops-snowflake-ext/master/docs/demo.png)

## Key Features

- **🔍 Live Search** - Real-time search with configurable debouncing
- **📊 Snowflake Integration** - Direct connection to your Snowflake data warehouse  
- **⚡ Fast Performance** - Optimized API calls with intelligent caching
- **⌨️ Keyboard Navigation** - Full keyboard support for accessibility
- **🎨 Modern UI** - Clean, responsive interface built with React

## Quick Setup

1. Install this extension
2. Add the "Snowflake Dropdown" control to your work item form
3. Configure your API endpoint and data type
4. Start searching your Snowflake data!

## Configuration Options

When adding the control to a work item form:

- **Field Name**: Work item field to bind selected values
- **API URL**: Your Go backend API endpoint (e.g., `https://api.yourcompany.com`)
- **Data Type**: Custom identifier for your data source (configure in your backend)
- **Search Settings**: Minimum search length (1-5 chars) and debounce delay (100-1000ms)

## Backend Implementation Required

This extension requires a backend API to connect to Snowflake. You have two options:

1. **Use the included Go backend** - Deploy the Go server from our repository
2. **Fork and customize** - Modify the project to use Snowflake SQL REST API directly

The included Go backend supports configurable data types that you define based on your Snowflake schema.

## API Response Format

Your backend API should return JSON in this format:

```json
{
  "results": [
    {
      "code": "ITEM001", 
      "description": "Your data description"
    }
  ]
}
```

The `code` and `description` fields can contain any data from your Snowflake tables.

## Get Started

Visit our [GitHub repository](https://github.com/jerrodtuck/azure-devops-snowflake-ext) for detailed setup instructions, API examples, and backend implementation guides.

---

**Need Help?** Check our documentation or open an issue on GitHub for support.