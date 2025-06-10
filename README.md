# qweather-mcp-go
[![smithery badge](https://smithery.ai/badge/@overstarry/qweather-mcp-go-hot)](https://smithery.ai/server/@overstarry/qweather-mcp-go-hot)

MCP (Model Context Protocol) tool implementation for QWeather API.

## Features

This project provides a series of QWeather API tools, including:

- Real-time weather query
- Weather forecast
- Air quality query
- Life indices query

## Running Methods

This project supports two running modes:

1. **stdio mode**: Communicate with clients through standard input and output
2. **SSE mode**: Provide API on HTTP server through Server-Sent Events (default mode)

### Environment Variables Setup

The following environment variables need to be set before running:

- `QWEATHER_API_BASE`: Base URL of QWeather API (e.g., `https://api.qweather.com`)
- `QWEATHER_API_KEY`: QWeather API key

### Windows Running Method

1. Edit the `run.bat` file to set your API key
2. Double-click to run the `run.bat` file or run `run.bat` in the command line

### Linux/Mac Running Method

1. Edit the `.env` file to set your API key
2. Run the following commands:

```bash
chmod +x run.sh
./run.sh
```

### Command Line Arguments

You can use the following command line arguments to control the program's behavior:

- `-t` or `--transport`: Specify the transport type, options are `stdio` or `sse` (default is `sse`)
- `-p` or `--port`: Specify the port for the SSE server to listen on (default is `8080`)

For example:

```bash
go run main.go -t stdio  # Run in stdio mode
go run main.go -p 3000   # Run in SSE mode, listening on port 3000
```

## Usage

### SSE Mode

When running in SSE mode, the server will provide HTTP API on the specified port. You can connect to this server using a client that supports the MCP protocol.

By default, the server address is: `http://localhost:8080`

### stdio Mode

When running in stdio mode, the server will communicate with clients through standard input and output. This mode is suitable for integration with AI assistants (such as Claude) that support the MCP protocol.