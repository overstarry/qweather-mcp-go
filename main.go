package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/overstarry/qweather-mcp-go/api"
	"github.com/overstarry/qweather-mcp-go/middlewares"
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
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "qweather",
		Version: "1.0.0",
	}, nil)

	// Register tools
	tools.RegisterWeatherTools(s, client)
	tools.RegisterAirQualityTools(s, client)
	tools.RegisterIndicesTools(s, client)

	// Start server based on transport type
	addr := ":" + port
	ctx := context.Background()
	switch transport {
	case "stdio":
		fmt.Println("QWeather MCP server running on stdio transport")
		if err := s.Run(ctx, &mcp.StdioTransport{}); err != nil {
			log.Fatal(err)
		}

	case "sse":
		baseURL := "http://localhost:" + port

		// Create SSE HTTP handler
		handler := mcp.NewSSEHandler(func(req *http.Request) *mcp.Server {
			return s
		}, &mcp.SSEOptions{})

		// Apply middlewares: recovery first, then logging
		handlerWithMiddlewares := middlewares.LoggingHandler(middlewares.RecoveryHandler(handler))
		fmt.Printf("QWeather MCP server running on SSE transport, listening at %s\n", addr)
		log.Printf("MCP server listening on %s", baseURL)
		fmt.Printf("SSE endpoint: %s\n", baseURL)

		// Create HTTP server with graceful shutdown support
		srv := &http.Server{
			Addr:    addr,
			Handler: handlerWithMiddlewares,
		}

		// Start server in background
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("SSE server error: %v", err)
			}
		}()

		// Wait for interrupt signal
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down SSE server...")

		// Graceful shutdown with 30 second timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatal("SSE server forced to shutdown:", err)
		}
		log.Println("SSE server exited")

	case "streamable":
		baseURL := "http://localhost:" + port

		// Create Streamable HTTP server (official implementation)
		handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
			return s
		}, nil)

		// Apply middlewares: recovery first, then logging
		handlerWithMiddlewares := middlewares.LoggingHandler(middlewares.RecoveryHandler(handler))
		fmt.Printf("QWeather MCP server running on Streamable HTTP transport, listening at %s\n", addr)
		log.Printf("MCP server listening on %s", baseURL)

		// Create HTTP server with graceful shutdown support
		srv := &http.Server{
			Addr:    addr,
			Handler: handlerWithMiddlewares,
		}

		// Start server in background
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Streamable server error: %v", err)
			}
		}()

		// Wait for interrupt signal
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down Streamable server...")

		// Graceful shutdown with 30 second timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatal("Streamable server forced to shutdown:", err)
		}
		log.Println("Streamable server exited")

	default:
		log.Fatalf("Unsupported transport type: %s", transport)
	}
}
