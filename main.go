package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/overstarry/qweather-mcp-go/api"
	"github.com/overstarry/qweather-mcp-go/tools"
)

func main() {
	// Command line arguments
	var transport string
	var port string
	flag.StringVar(&transport, "t", "sse", "Transport type (stdio or sse)")
	flag.StringVar(&transport, "transport", "sse", "Transport type (stdio or sse)")
	flag.StringVar(&port, "p", "8080", "SSE server listening port")
	flag.StringVar(&port, "port", "8080", "SSE server listening port")
	flag.Parse()

	// Get configuration from environment variables
	baseURL := os.Getenv("QWEATHER_API_BASE")
	apiKey := os.Getenv("QWEATHER_API_KEY")

	if baseURL == "" || apiKey == "" {
		log.Fatal("Environment variables QWEATHER_API_BASE and QWEATHER_API_KEY must be set")
	}

	// Create API client
	client := api.NewClient(baseURL, apiKey)

	// Create MCP server
	s := server.NewMCPServer(
		"qweather",
		"1.0.0",
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Register tools
	tools.RegisterWeatherTools(s, client)
	tools.RegisterAirQualityTools(s, client)
	tools.RegisterIndicesTools(s, client)

	// Start server based on transport type
	if transport == "stdio" {
		fmt.Println("QWeather MCP server running on stdio")
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		// Default to SSE
		addr := ":" + port
		baseURL := "http://localhost:" + port
		sseServer := server.NewSSEServer(s, server.WithBaseURL(baseURL))
		fmt.Printf("QWeather MCP server running on SSE, listening at %s\n", addr)
		if err := sseServer.Start(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}
