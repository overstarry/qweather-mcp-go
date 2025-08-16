# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

qweather-mcp-go is an MCP (Model Context Protocol) server implementation that provides weather data tools using the QWeather API. It serves as a bridge between AI assistants and weather services, offering real-time weather, forecasts, air quality, and life indices data.

## Development Commands

### Building and Running
- `go run main.go` - Run in SSE mode (default, port 8080)
- `go run main.go -t stdio` - Run in stdio mode for MCP integration
- `go run main.go -t sse` - Run in SSE mode (default)
- `go run main.go -t streamable` - Run in Streamable HTTP mode
- `go run main.go -p 3000` - Run on port 3000 (for sse and streamable modes)
- `go run main.go -t streamable -p 8080` - Run Streamable HTTP on port 8080
- `go build -o qweather-mcp-go main.go` - Build binary

### Environment Setup
- Copy `.env.example` to `.env` and set:
  - `QWEATHER_API_BASE=https://api.qweather.com`
  - `QWEATHER_API_KEY=your_api_key_here`
- Or set environment variables directly:
  ```bash
  export QWEATHER_API_BASE=https://api.qweather.com
  export QWEATHER_API_KEY=your_api_key_here
  ```

### Testing
- `go test ./...` - Run all tests
- `go test ./utils -v` - Run utils tests with verbose output
- `go test -race ./...` - Run tests with race detection

### Code Quality
- `go fmt ./...` - Format all Go files
- `go vet ./...` - Run Go vet for static analysis
- `go mod tidy` - Clean up module dependencies

### Platform-specific Scripts
- Linux/Mac: `./run.sh` (loads .env automatically)
- Windows: `run.bat` (loads .env automatically)

### Docker
- `docker build -t qweather-mcp-go .` - Build Docker image
- Uses multi-stage build with golang:1.23-alpine and debian:bookworm-slim

## Architecture

### Core Components

**api/** - QWeather API client and data models
- `client.go` - HTTP client with logging, handles all QWeather API endpoints
- `models.go` - Complete Go structs for all API response formats
- `logger.go` - Logging levels (Error, Info, Debug)

**tools/** - MCP tool implementations
- `weather.go` - Weather tools (current, forecast, hourly, minutely, warnings)
- `air_quality.go` - Air quality tools (current, hourly, daily forecasts)
- `indices.go` - Life indices tools (UV, comfort, etc.)

**utils/** - Utility functions
- `helpers.go` - String joining utilities
- `helpers_test.go` - Unit tests for utilities

### Transport Modes
- **stdio mode**: For MCP integration with AI assistants
- **SSE mode**: HTTP server with Server-Sent Events (default)
- **streamable HTTP mode**: HTTP server with streamable responses for web applications and REST-like APIs

### Tool Categories
1. **Weather Tools**: Real-time conditions, daily/hourly forecasts, minutely precipitation, weather warnings
2. **Air Quality Tools**: Current air quality, hourly/daily AQI forecasts with pollutant details
3. **Life Indices Tools**: UV index, comfort level, clothing suggestions, etc.

### Error Handling
- All API responses include status codes and error messages
- Client includes configurable logging levels
- Tools return proper MCP error responses for invalid inputs

### Configuration
- Environment variables for API credentials
- Command-line flags for transport mode (`stdio`, `sse`, `streamable`) and port  
- Smithery integration for easy deployment via `npx -y @smithery/cli install @overstarry/qweather-mcp-go --client claude`

### Transport Endpoints
- **Stdio**: Standard input/output communication
- **SSE**: `http://localhost:PORT/sse` - Server-Sent Events endpoint
- **Streamable HTTP**: `http://localhost:PORT` - Streamable HTTP endpoint with native MCP over HTTP support
  - Uses the official `NewStreamableHTTPServer` implementation from mcp-go v0.37.0

### Key Implementation Details
- MCP server uses `github.com/mark3labs/mcp-go v0.37.0` framework
- API client includes configurable timeout (10s default) and logging levels
- All tools validate input parameters and return structured MCP responses
- City lookups use QWeather's geocoding API before querying weather data
- Transport modes: stdio (for MCP clients), SSE (HTTP server for real-time updates), and streamable HTTP (REST-like API server)