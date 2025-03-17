package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sergeymakinen/go-ico"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

func main() {
	// Get the appropriate log directory for the current OS
	logDir, err := os.UserConfigDir()
	if err != nil {
		logDir = os.TempDir()
	}
	logDir = filepath.Join(logDir, "favicon-mcp-server")

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatal("Failed to create log directory:", err)
	}

	// Create log file
	logFile, err := os.OpenFile(filepath.Join(logDir, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()

	// Redirect log output to file
	log.SetOutput(logFile)

	s := server.NewMCPServer("FaviconGenerator", "1.0.0")

	tool := mcp.NewTool(
		"svg_to_favicon",
		mcp.WithDescription("Converts SVG icons into various standard website favicon formats such as ICO and PNG."),
		mcp.WithString("svg_data", mcp.Description("SVG icon content provided as a string.")),
		mcp.WithString("svg_file", mcp.Description("Path to the SVG file to convert.")),
		mcp.WithString("output_dir", mcp.Description("Directory to save the output files. If not provided, returns base64 encoded data.")),
		mcp.WithArray("output_formats", mcp.Description("An array of strings specifying the desired output formats.")),
	)
	s.AddTool(tool, faviconHandler)

	log.Println("Starting MCP server...")
	if err := server.ServeStdio(s); err != nil {
		log.Printf("Server error: %v\n", err)
	}
}

// FaviconHandler handles the invocation of the svg_to_favicon tool
func faviconHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get SVG data either from direct input or file
	var svgData string
	if svgPath, ok := request.Params.Arguments["svg_file"].(string); ok && svgPath != "" {
		// Read SVG from file
		log.Println("Reading SVG from file:", svgPath)
		data, err := os.ReadFile(svgPath)
		if err != nil {
			log.Println("Error reading SVG file:", err)
			return mcp.NewToolResultError(fmt.Sprintf("Failed to read SVG file: %v", err)), nil
		}
		svgData = string(data)
		log.Println("Successfully read SVG file, content length:", len(svgData))
	} else {
		// Get SVG from direct input
		svgData, ok := request.Params.Arguments["svg_data"].(string)
		if !ok || svgData == "" {
			log.Println("No SVG data provided")
			return mcp.NewToolResultError("Either svg_data or svg_file parameter is required"), nil
		}
		log.Println("Using direct SVG input, content length:", len(svgData))
	}

	// Define standard favicon sizes and formats
	type faviconFormat struct {
		name string
		size int
	}

	// Standard favicon formats
	formats := []faviconFormat{
		{"android-chrome", 192},
		{"android-chrome", 512},
		{"apple-touch-icon", 180},
		{"favicon", 16},
		{"favicon", 32},
	}

	// Map to store PNG images
	pngImages := make(map[string]*bytes.Buffer)
	// Images for ICO format (16x16 and 32x32)
	var icoImages []image.Image

	// Check if SVG content is valid
	if !strings.Contains(svgData, "<svg") {
		log.Println("Error: Invalid SVG content - missing <svg> tag")
		return mcp.NewToolResultError("Invalid SVG content"), nil
	}

	// Print SVG content preview
	previewLen := 200
	if len(svgData) > previewLen {
		log.Printf("SVG content preview (first %d chars): %s...\n", previewLen, svgData[:previewLen])
	} else {
		log.Printf("SVG content (full, %d chars): %s\n", len(svgData), svgData)
	}

	// Parse SVG and create icon
	icon, err := oksvg.ReadIconStream(strings.NewReader(svgData))
	if err != nil {
		log.Println("Failed to parse SVG:", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to parse SVG: %v", err)), nil
	}

	// Check for SVG elements and their attributes
	elements := map[string][]string{
		"circle":         {"cx", "cy", "r", "fill", "stroke"},
		"text":           {"x", "y", "font-family", "font-size", "fill"},
		"linearGradient": {"id", "x1", "y1", "x2", "y2"},
		"stop":           {"offset", "stop-color", "stop-opacity"},
	}

	// Check SVG content
	hasText := false
	hasGradient := false

	for elem, attrs := range elements {
		if strings.Contains(svgData, "<"+elem) {
			log.Printf("Found SVG element: <%s>\n", elem)
			for _, attr := range attrs {
				if strings.Contains(svgData, attr) {
					log.Printf("  - Has attribute: %s\n", attr)
				}
			}
			if elem == "text" {
				hasText = true
			} else if elem == "linearGradient" {
				hasGradient = true
			}
		}
	}

	// Print warnings for potential issues
	if hasText {
		log.Println("Warning: SVG contains text elements which may not render correctly in small sizes")
	}
	if hasGradient {
		log.Println("Warning: SVG contains gradients which may not render correctly in all browsers")
	}

	// Print icon information
	log.Printf("SVG dimensions: ViewBox=%.2fx%.2f, Width=%.2f, Height=%.2f\n",
		icon.ViewBox.W, icon.ViewBox.H, icon.ViewBox.W, icon.ViewBox.H)

	// Get original dimensions and set icon parameters
	origW, origH := icon.ViewBox.W, icon.ViewBox.H
	if origW == 0 || origH == 0 {
		log.Println("Error: Invalid SVG dimensions - width or height is 0")
		return mcp.NewToolResultError("Invalid SVG dimensions"), nil
	}

	// Generate different sizes directly from SVG
	for _, format := range formats {
		// Create a new RGBA image
		img := image.NewRGBA(image.Rect(0, 0, format.size, format.size))

		// Calculate scale to fit while maintaining aspect ratio
		scaleX := float64(format.size) / origW
		scaleY := float64(format.size) / origH
		scale := scaleX
		if scaleY < scaleX {
			scale = scaleY
		}

		// Calculate centered position
		x := (float64(format.size) - (origW * scale)) / 2
		y := (float64(format.size) - (origH * scale)) / 2

		// Create a new icon instance for each size
		sizeIcon, err := oksvg.ReadIconStream(strings.NewReader(svgData))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to parse SVG: %v", err)), nil
		}

		// Create rasterizer with antialiasing
		scanner := rasterx.NewScannerGV(format.size, format.size, img, img.Bounds())
		scanner.SetClip(img.Bounds())
		r := rasterx.NewDasher(format.size, format.size, scanner)

		// Set icon size and position
		sizeIcon.SetTarget(x, y, origW*scale, origH*scale)

		// Draw icon with high quality settings
		sizeIcon.Draw(r, 1.0)

		// Encode PNG
		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create PNG image: %v", err)), nil
		}

		// Save to pngImages with format name
		fileName := fmt.Sprintf("%s-%dx%d.png", format.name, format.size, format.size)
		if format.name == "apple-touch-icon" {
			fileName = "apple-touch-icon.png"
		}
		pngImages[fileName] = buf

		// For ICO format, collect 16x16 and 32x32 images
		if format.size == 16 || format.size == 32 {
			icoImages = append(icoImages, img)
		}
	}

	results := make(map[string]interface{})

	// Get output directory
	outputDir, ok := request.Params.Arguments["output_dir"].(string)
	log.Println("Output directory:", outputDir, "(exists:", ok, ")")
	if ok && outputDir != "" {
		log.Println("Creating output directory:", outputDir)
		// Create output directory if it doesn't exist
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Println("Error creating output directory:", err)
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create output directory: %v", err)), nil
		}
		log.Println("Successfully created output directory")

		// Save PNG files
		for fileName, buf := range pngImages {
			if err := os.WriteFile(filepath.Join(outputDir, fileName), buf.Bytes(), 0644); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to write PNG file: %v", err)), nil
			}
		}

		// Save ICO file
		// Save ICO file if we have the required images
		if len(icoImages) > 0 {
			icoBuf := new(bytes.Buffer)
			if err := ico.EncodeAll(icoBuf, icoImages); err != nil {
				log.Println("ICO encoding error:", err)
				results["ico_error"] = err.Error()
			} else {
				if err := os.WriteFile(filepath.Join(outputDir, "favicon.ico"), icoBuf.Bytes(), 0644); err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("Failed to write ICO file: %v", err)), nil
				}
				results["ico"] = base64.StdEncoding.EncodeToString(icoBuf.Bytes())
			}
		}
	} else {
		log.Println("No output directory specified, will return base64 encoded data")
		// Return base64 encoded data
		pngResult := make(map[string]string)
		for fileName, buf := range pngImages {
			pngResult[fileName] = base64.StdEncoding.EncodeToString(buf.Bytes())
		}
		results["png"] = pngResult

		// Generate ICO file if we have the required images
		if len(icoImages) > 0 {
			icoBuf := new(bytes.Buffer)
			if err := ico.EncodeAll(icoBuf, icoImages); err != nil {
				log.Println("ICO encoding error:", err)
				results["ico_error"] = err.Error()
			} else {
				results["ico"] = base64.StdEncoding.EncodeToString(icoBuf.Bytes())
			}
		}
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
