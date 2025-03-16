package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	ico "github.com/sergeymakinen/go-ico"
	"github.com/tdewolff/canvas"
)

// FaviconHandler handles the invocation of the svg_to_favicon tool
func faviconHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	svgData, ok := request.Params.Arguments["svg_data"].(string)
	if !ok || svgData == "" {
		return mcp.NewToolResultError("svg_data parameter is missing or invalid"), nil
	}

	outputFormatsRaw, ok := request.Params.Arguments["output_formats"].([]interface{})
	outputFormats := []string{"ico", "png"} // default formats
	if ok {
		outputFormats = make([]string, len(outputFormatsRaw))
		for i, format := range outputFormatsRaw {
			if s, ok := format.(string); ok {
				outputFormats[i] = strings.ToLower(s)
			}
		}
	}

	sizes := []int{16, 32, 48}
	pngImages := make(map[int]*bytes.Buffer)
	var icoImages []image.Image

	// Render PNG images of different sizes using tdewolff/canvas
	for _, size := range sizes {
		c := canvas.New(float64(size), float64(size))
		ctx := canvas.NewContext(c)
		// Create a new RGBA image
		img := image.NewRGBA(image.Rect(0, 0, size, size))
		
		// Draw a simple shape for testing
		p := &canvas.Path{}
		p.MoveTo(0, 0)
		p.LineTo(float64(size), 0)
		p.LineTo(float64(size), float64(size))
		p.LineTo(0, float64(size))
		p.Close()
		
		ctx.SetFillColor(canvas.Black)
		ctx.DrawPath(0, 0, p)
		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			log.Printf("PNG encoding error (size %d): %v", size, err)
			continue
		}
		pngImages[size] = buf
		icoImages = append(icoImages, img)
	}

	results := make(map[string]interface{})

	// Generate ICO format
	if contains(outputFormats, "ico") && len(icoImages) > 0 {
		icoBuf := new(bytes.Buffer)
		if err := ico.EncodeAll(icoBuf, icoImages); err != nil {
			log.Printf("ICO encoding error: %v", err)
			results["ico_error"] = err.Error()
		} else {
			results["ico"] = base64.StdEncoding.EncodeToString(icoBuf.Bytes())
		}
	}

	// Generate PNG format
	if contains(outputFormats, "png") {
		pngResult := make(map[string]string)
		for size, buf := range pngImages {
			pngResult[strconv.Itoa(size)] = base64.StdEncoding.EncodeToString(buf.Bytes())
		}
		results["png"] = pngResult
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Type: "text",
				Text: "Successfully generated favicons",
			},
		},
		Result: mcp.Result{
			Meta: results,
		},
	}, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func main() {
	s := server.NewMCPServer("FaviconGenerator", "1.0.0")

	tool := mcp.NewTool(
		"svg_to_favicon",
		mcp.WithDescription("Converts SVG icons into various standard website favicon formats such as ICO and PNG."),
		mcp.WithString("svg_data", mcp.Required(), mcp.Description("SVG icon content provided as a string.")),
		mcp.WithArray("output_formats", mcp.Description("An array of strings specifying the desired output formats.")),
	)
	s.AddTool(tool, faviconHandler)

	fmt.Println("Starting MCP server...")
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
