# Favicon MCP Server

A Model Control Protocol (MCP) server that converts SVG images into various favicon formats (ICO and PNG) for web applications.

## Features

- **SVG to ICO Conversion**: Convert SVG images to ICO format (16x16, 32x32, 48x48 pixels).
- **SVG to PNG Conversion**: Convert SVG images to PNG format (16x16, 32x32, 48x48 pixels).
- **Base64 Encoded Output**: Easy integration with base64 encoded output.
- **MCP Protocol Support**: Seamless integration with LLM-powered applications.


## Prerequisites

- **Go 1.20 or higher**
- **Dependencies**:
    - `github.com/mark3labs/mcp-go v0.13.0`
    - `github.com/sergeymakinen/go-ico`
    - `github.com/tdewolff/canvas`


## Installation

```bash
git clone https://github.com/elliotxx/favicon-mcp-server.git
cd favicon-mcp-server
go mod download
```


## Usage

- **Start the server**:

```bash
go run main.go
```

- **Tool Parameters**:
    - `svg_data`: SVG icon content provided as a string.
    - `output_formats`: Array of strings specifying the desired output formats (`ico`, `png`). Default: `["ico", "png"]`.


### Example Response

The server returns base64 encoded favicon data in the requested formats:

```json
{
  "content": [
    {
      "type": "text",
      "text": "Successfully generated favicons"
    }
  ],
  "meta": {
    "ico": "base64_encoded_ico_data",
    "png": "base64_encoded_png_data"
  }
}
```

## Usage

### Input Methods

1. Direct SVG Input:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "svg_to_favicon",
    "arguments": {
      "svg_data": "<svg width=\"32\" height=\"32\"><rect width=\"32\" height=\"32\" fill=\"red\"/></svg>"
    }
  }
}
```

2. SVG File Input:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "svg_to_favicon",
    "arguments": {
      "svg_file": "path/to/your/icon.svg"
    }
  }
}
```

### Output Methods

1. Base64 Encoded Output (Default):
   - Returns base64 encoded ICO and PNG data in the response

2. File Output:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "svg_to_favicon",
    "arguments": {
      "svg_data": "<svg width=\"32\" height=\"32\"><rect width=\"32\" height=\"32\" fill=\"red\"/></svg>",
      "output_dir": "path/to/output/directory",
      "output_formats": ["ico", "png"]
    }
  }
}
```

When using file output:
- ICO file will be saved as `favicon.ico`
- PNG files will be saved as `favicon-{size}x{size}.png` (e.g., `favicon-32x32.png`)

## Testing

### Quick Test (Recommended)

The simplest way to test with direct SVG input:
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"svg_to_favicon","arguments":{"svg_data":"<svg width=\"32\" height=\"32\"><rect width=\"32\" height=\"32\" fill=\"red\"/></svg>"}}}' | go run main.go
```

### Using Test File

If you prefer using a test file:

1. Create a test file `test.json` with your test case
2. Run the command:
```bash
echo $(tr -d '\n' < test/test.json) | go run main.go
```

### Parameters

- `svg_data`: SVG content as a string
- `svg_file`: Path to an SVG file
- `output_dir`: Directory to save the output files
- `output_formats`: Array of desired formats (`["ico", "png"]`)

## Development

The project follows standard Go project layout and uses Go modules for dependency management.

### Project Structure

```
favicon-mcp-server/
├── main.go         # Main server implementation
├── go.mod         # Go module definition
├── go.sum         # Go module checksums
└── README.md      # This file
```

### Building from Source

```bash
go build
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request


### Integration with Windsurf

To integrate this MCP server with Windsurf, follow these steps:

1. **Open Windsurf** and navigate to the Cascade interface.
2. **Configure MCP Servers**:
    - Open the `~/.codeium/windsurf/mcp_config.json` file by clicking the hammer icon and selecting "Configure".
    - Add the following configuration:

```json
"mcpServers": {
  "favicon-mcp-server": {
    "command": "go",
    "args": ["run", "main.go"],
    "cwd": "/path/to/favicon-mcp-server",
    "env": {}
  }
}
```

Replace `/path/to/favicon-mcp-server` with the actual path to your project directory.
3. **Refresh Windsurf**:
    - Click the "Refresh" button in the MCP toolbar to load the new configuration.

### Integration with Cursor

To integrate this MCP server with Cursor, follow these steps:

1. **Enable MCP Servers**:
    - Navigate to Cursor settings and find the MCP servers option.
    - Enable MCP servers if not already enabled.
2. **Add New MCP Server**:
    - Click "Add new MCP server".
    - Provide the path to your executable or command to run the server.
    - For this project, you might need to create a standalone executable or use a bundling tool to simplify integration.
3. **Configure Server Details**:
    - Enter the command to run your MCP server. For example:

```bash
go run main.go
```

    - Ensure the path to the executable is correct.
4. **Enable the Server**:
    - After adding the server, click "Enable" to activate it.

By following these steps, you can integrate the Favicon MCP Server with both Windsurf and Cursor, enhancing your development workflow with AI-powered tools.

## License

This project is licensed under the MIT License - see the LICENSE file for details