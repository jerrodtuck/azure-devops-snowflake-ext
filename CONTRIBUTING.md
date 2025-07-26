# Contributing to Azure DevOps Dynamic Dropdown

First off, thank you for considering contributing! 🎉

## How Can I Contribute?

### 🐛 Reporting Bugs

Before creating bug reports, please check existing issues. When creating a bug report, include:

- Clear descriptive title
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, database type, versions)
- Logs/screenshots if applicable

### 💡 Suggesting Features

Feature requests are welcome! Please provide:

- Clear use case
- Why this would be useful for others
- Possible implementation approach

### 🔧 Pull Requests

1. Fork the repo
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

#### PR Guidelines

- Update README.md with details of changes if needed
- Add tests for new functionality
- Ensure all tests pass
- Update documentation

### 📦 Adding Database Connectors

To add support for a new database:

1. Create connector in `connectors/yourdb/connector.go`
2. Implement the `DataSource` interface:
   ```go
   type DataSource interface {
       Connect() error
       Query(sql string) ([]DropdownItem, error)
       Close() error
   }
   ```
3. Add configuration in `.env.example`
4. Update README with setup instructions
5. Add tests

Example connector structure:
```
connectors/
  └── mongodb/
      ├── connector.go
      ├── connector_test.go
      └── README.md
```

### 🌐 Adding Translations

1. Add language file in `extension/locales/`
2. Update supported languages in README

## Development Setup

```bash
# Clone your fork
git clone https://github.com/your-username/azure-devops-dropdown.git
cd azure-devops-dropdown

# Install dependencies
go mod download
cd extension && npm install

# Run tests
go test ./...
npm test

# Run locally
docker-compose -f docker-compose.dev.yml up
```

## Code Style

- Go: Use `gofmt` and `golint`
- JavaScript: Use provided ESLint config
- Commit messages: Use conventional commits

## Testing

- Unit tests for all new functions
- Integration tests for database connectors
- Manual testing with Azure DevOps

## Documentation

- Code comments for complex logic
- README updates for new features
- API documentation for new endpoints

## Questions?

Feel free to open an issue with the "question" label or start a discussion!

Thank you! 🙏