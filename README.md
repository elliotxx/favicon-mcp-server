# Favicon MCP Server

A Model Control Protocol (MCP) server that converts SVG images into various favicon formats (ICO and PNG) for web applications.

## Features

- Convert SVG images to ICO format (16x16, 32x32, 48x48 pixels)
- Convert SVG images to PNG format (16x16, 32x32, 48x48 pixels)
- Base64 encoded output for easy integration
- MCP protocol support for seamless integration with LLM-powered applications

## Prerequisites

- Go 1.20 or higher
- Dependencies:
  - github.com/mark3labs/mcp-go v0.13.0
  - github.com/sergeymakinen/go-ico
  - github.com/tdewolff/canvas

## Installation

```bash
git clone https://github.com/yourusername/favicon-mcp-server.git
cd favicon-mcp-server
go mod download
```

## Usage

1. Start the server:
```bash
go run main.go
```

2. The server provides one tool: `svg_to_favicon`

### Tool Parameters

- `svg_data` (required): SVG icon content provided as a string
- `output_formats` (optional): Array of strings specifying the desired output formats
  - Supported formats: `ico`, `png`
  - Default: `["ico", "png"]`

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

## Testing

### Quick Test (Recommended)

The simplest and most reliable way to test:
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"svg_to_favicon","arguments":{"svg_data":"<svg width=\"32\" height=\"32\"><rect width=\"32\" height=\"32\" fill=\"red\"/></svg>"}}}' | go run main.go
```

### Using Test File

If you prefer using a test file:

1. Create a test file `test.json` with your test case, for example:
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

2. Test using pipe (make sure to remove newlines):
```bash
tr -d '\n' < test.json | go run main.go
```

### Optional Parameters

You can also specify output formats in the arguments:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "svg_to_favicon",
    "arguments": {
      "svg_data": "<svg width=\"32\" height=\"32\"><rect width=\"32\" height=\"32\" fill=\"red\"/></svg>",
      "output_formats": ["ico", "png"]
    }
  }
}
```

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

## License

This project is licensed under the MIT License - see the LICENSE file for details.
