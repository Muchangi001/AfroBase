package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type ImagePayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func main() {
	// Create Fiber instance
	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024, // 50MB limit for large images
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Create uploads directory if it doesn't exist
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Fatal("Failed to create uploads directory:", err)
	}

	// Upload endpoint
	app.Post("/upload", handleImageUpload)

	// Health check endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Fiber Image Server is running",
			"port":    "5175",
		})
	})

	// API endpoint to get image list
	app.Get("/api/images", getImageList)

	// Serve static files from uploads directory
	app.Static("/uploads", "./uploads")

	// Start server
	log.Println("Server starting on port 5175...")
	log.Fatal(app.Listen(":5174"))
}

func getImageList(c *fiber.Ctx) error {
	// read all files in the uploads directory
	files, err := ioutil.ReadDir("./uploads")
	if err != nil {
		log.Printf("Error reading uploads directory: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to read uploads directory",
			"success": false,
		})
	}

	// send images in uploads directory as JSON
	var images []map[string]interface{} = make([]map[string]interface{}, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			// Get file info
			fileInfo, err := os.Stat(filepath.Join("./uploads", file.Name()))
			if err != nil {
				log.Printf("Error getting file info: %v", err)
				continue
			}

			// Create image object
			image := map[string]interface{}{
				"name":        file.Name(),
				"size":        fileInfo.Size(),
				"upload_time": fileInfo.ModTime().Unix(),
				"title":       strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())),
				"description": "Uploaded image",
				"url":         "http://localhost:5174/uploads/" + file.Name(),
			}
			images = append(images, image)
		}
	}

	// Return images as JSON
	if len(images) == 0 {
		return c.JSON([]map[string]interface{}{})
	}
	return c.JSON(images)
}

func handleImageUpload(c *fiber.Ctx) error {
	var payload ImagePayload

	// Parse JSON body
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	// Validate payload
	if payload.Image == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Image data is required",
			"success": false,
		})
	}

	// Decode base64 image
	imageData, err := base64.StdEncoding.DecodeString(payload.Image)
	if err != nil {
		log.Printf("Error decoding base64 image: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid base64 image data",
			"success": false,
		})
	}

	// Detect image format from first few bytes
	var fileExt string
	if len(imageData) >= 4 {
		switch {
		case imageData[0] == 0xFF && imageData[1] == 0xD8:
			fileExt = ".jpg"
		case imageData[0] == 0x89 && imageData[1] == 0x50 && imageData[2] == 0x4E && imageData[3] == 0x47:
			fileExt = ".png"
		case imageData[0] == 0x47 && imageData[1] == 0x49 && imageData[2] == 0x46:
			fileExt = ".gif"
		case imageData[0] == 0x52 && imageData[1] == 0x49 && imageData[2] == 0x46 && imageData[3] == 0x46:
			fileExt = ".webp"
		default:
			fileExt = ".jpg" // Default fallback
		}
	} else {
		fileExt = ".jpg"
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	sanitizedTitle := sanitizeFilename(payload.Title)
	if sanitizedTitle == "" {
		sanitizedTitle = "image"
	}
	filename := fmt.Sprintf("%d_%s%s", timestamp, sanitizedTitle, fileExt)
	filepath := filepath.Join("./uploads", filename)

	// Save file
	if err := ioutil.WriteFile(filepath, imageData, 0644); err != nil {
		log.Printf("Error saving file: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to save image",
			"success": false,
		})
	}

	// Log successful upload
	log.Printf("Image uploaded successfully: %s (Title: %s, Description: %s)", 
		filename, payload.Title, payload.Description)

	// Return success response
	return c.JSON(fiber.Map{
		"success": true,
		"url":     "/uploads/" + filename,
	})
}

// sanitizeFilename removes or replaces invalid characters for filenames
func sanitizeFilename(filename string) string {
	// Remove or replace invalid characters
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "/", "-")
	filename = strings.ReplaceAll(filename, "\\", "-")
	filename = strings.ReplaceAll(filename, ":", "-")
	filename = strings.ReplaceAll(filename, "*", "-")
	filename = strings.ReplaceAll(filename, "?", "-")
	filename = strings.ReplaceAll(filename, "\"", "-")
	filename = strings.ReplaceAll(filename, "<", "-")
	filename = strings.ReplaceAll(filename, ">", "-")
	filename = strings.ReplaceAll(filename, "|", "-")
	
	// Limit length
	if len(filename) > 50 {
		filename = filename[:50]
	}
	
	return filename
}
