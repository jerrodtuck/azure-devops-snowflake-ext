# PowerShell script to run the Go backend with .env loading
$env:PATH += ";C:\Program Files\Go\bin"

# Check if Go is installed
if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Please install Go from https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}

# Navigate to the Go backend directory
Set-Location $PSScriptRoot

Write-Host "Starting Snowflake Go Backend..." -ForegroundColor Green

# Run the Go application with all source files
go run main.go loadenv.go