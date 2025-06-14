#!/bin/bash

# Load environment variables
export $(grep -v '^#' .env | xargs)

# Display loaded environment variables
echo "Loaded environment variables:"
echo "QWEATHER_API_BASE=$QWEATHER_API_BASE"
echo "QWEATHER_API_KEY=$QWEATHER_API_KEY"

# Run the program
echo "Starting QWeather MCP server..."
go run main.go
