package cloudinary

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Client struct {
	cld *cloudinary.Cloudinary
}

// boolPtr returns a pointer to the given bool.
func boolPtr(b bool) *bool { return &b }

func NewClient() (*Client, error) {
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		cloudinaryURL = "cloudinary://436377263181225:y53nJ5-e01HR93f2CmcFHKO3lb4@dll1pjjbm"
	}

	fmt.Printf("🔍 Initializing Cloudinary client...\n")

	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloudinary client: %w", err)
	}

	fmt.Printf("✅ Cloudinary client initialized successfully\n")
	return &Client{cld: cld}, nil
}

type UploadResult struct {
	URL      string
	PublicID string
	Error    error
	Index    int
}

// cleanBase64Data removes the data:image/...;base64, prefix if present
func cleanBase64Data(imageData string) (string, error) {
	// If it contains a base64 prefix, extract just the base64 data
	if strings.Contains(imageData, ";base64,") {
		parts := strings.Split(imageData, ";base64,")
		if len(parts) == 2 {
			// Validate that it's valid base64
			_, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				return "", fmt.Errorf("invalid base64 data: %w", err)
			}
			return parts[1], nil
		}
	}

	// If it's already base64, validate it
	if len(imageData) > 0 {
		_, err := base64.StdEncoding.DecodeString(imageData)
		if err != nil {
			// Not valid base64, might be a URL
			if strings.HasPrefix(imageData, "http://") || strings.HasPrefix(imageData, "https://") {
				return imageData, nil
			}
			return "", fmt.Errorf("invalid image data format")
		}
	}

	return imageData, nil
}

func (c *Client) UploadImage(ctx context.Context, imageData string, folder string) (string, string, error) {
	// Clean the base64 data
	cleanData, err := cleanBase64Data(imageData)
	if err != nil {
		return "", "", fmt.Errorf("failed to clean image data: %w", err)
	}

	// Validate the data
	if len(cleanData) == 0 {
		return "", "", fmt.Errorf("empty image data after cleaning")
	}

	if len(cleanData) < 10 {
		return "", "", fmt.Errorf("image data too short (length: %d)", len(cleanData))
	}

	fmt.Printf("📤 Uploading image to Cloudinary (data length: %d chars)...\n", len(cleanData))

	// Create upload parameters
	params := uploader.UploadParams{
		Folder:         folder,
		Eager:          "f_auto,q_auto",
		UseFilename:    boolPtr(false), // Don't use filename from the string
		UniqueFilename: boolPtr(true),  // Generate unique filename
		ResourceType:   "image",
	}

	uploadCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var result *uploader.UploadResult

	// If it's a URL, upload as URL
	if strings.HasPrefix(cleanData, "http://") || strings.HasPrefix(cleanData, "https://") {
		result, err = c.cld.Upload.Upload(uploadCtx, cleanData, params)
	} else {
		// It's base64, use the data:image prefix
		result, err = c.cld.Upload.Upload(uploadCtx, "data:image/jpeg;base64,"+cleanData, params)
	}

	if err != nil {
		fmt.Printf("❌ Cloudinary upload error: %v\n", err)
		return "", "", fmt.Errorf("cloudinary upload failed: %w", err)
	}

	fmt.Printf("✅ Cloudinary upload success: URL=%s, PublicID=%s\n", result.SecureURL, result.PublicID)
	return result.SecureURL, result.PublicID, nil
}

func (c *Client) UploadMultipleImages(ctx context.Context, images []string, folder string) ([]UploadResult, error) {
	if len(images) == 0 {
		return nil, nil
	}

	fmt.Printf("📤 Starting upload of %d images to Cloudinary...\n", len(images))

	resultChan := make(chan UploadResult, len(images))
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// Limit concurrent uploads to avoid rate limiting
	semaphore := make(chan struct{}, 3)

	for i, imageData := range images {
		if strings.TrimSpace(imageData) == "" {
			fmt.Printf("⚠️ Image %d is empty, skipping\n", i+1)
			continue
		}

		wg.Add(1)
		go func(index int, data string) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			fmt.Printf("📤 Uploading image %d to Cloudinary...\n", index+1)

			url, publicID, err := c.UploadImage(ctx, data, folder)

			if err != nil {
				fmt.Printf("❌ Image %d upload failed: %v\n", index+1, err)
				resultChan <- UploadResult{Index: index, Error: err}
				return
			}

			fmt.Printf("✅ Image %d uploaded: URL=%s\n", index+1, url)
			resultChan <- UploadResult{
				Index:    index,
				URL:      url,
				PublicID: publicID,
				Error:    nil,
			}
		}(i, imageData)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []UploadResult
	for result := range resultChan {
		results = append(results, result)
	}

	// Count successful uploads
	var successCount, failCount int
	for _, r := range results {
		if r.Error != nil {
			failCount++
		} else {
			successCount++
		}
	}

	fmt.Printf("📊 Upload complete: %d successful, %d failed\n", successCount, failCount)

	return results, nil
}
