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
	flag.StringVar(&transport, "t", "sse", "Transport type (stdio, sse, or streamable)")
	flag.StringVar(&transport, "transport", "sse", "Transport type (stdio, sse, or streamable)")
	flag.StringVar(&port, "p", "8080", "Server listening port (for sse and streamable transports)")
	flag.StringVar(&port, "port", "8080", "Server listening port (for sse and streamable transports)")
	flag.Parse()

	// Validate transport type
	if transport != "stdio" && transport != "sse" && transport != "streamable" {
		log.Fatalf("Invalid transport type: %s. Must be one of: stdio, sse, streamable", transport)
	}

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
	addr := ":" + port

	switch transport {
	case "stdio":
		fmt.Println("QWeather MCP server running on stdio transport")
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Stdio server error: %v", err)
		}

	case "sse":
		baseURL := "http://localhost:" + port
		sseServer := server.NewSSEServer(s, server.WithBaseURL(baseURL))
		fmt.Printf("QWeather MCP server running on SSE transport, listening at %s\n", addr)
		fmt.Printf("SSE endpoint: %s/sse\n", baseURL)
		if err := sseServer.Start(addr); err != nil {
			log.Fatalf("SSE server error: %v", err)
		}

	case "streamable":
		baseURL := "http://localhost:" + port
		fmt.Printf("QWeather MCP server running on Streamable HTTP transport, listening at %s\n", addr)
		fmt.Printf("HTTP endpoint: %s\n", baseURL)
		
		// Create Streamable HTTP server (official implementation)
		httpServer := server.NewStreamableHTTPServer(s)
		if err := httpServer.Start(addr); err != nil {
			log.Fatalf("Streamable HTTP server error: %v", err)
		}

	default:
		log.Fatalf("Unsupported transport type: %s", transport)
	}
}
