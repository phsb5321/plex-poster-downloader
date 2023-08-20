// Package downloader provides functionality to download and save images from provided URLs.
package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadImage downloads an image from the provided URL and saves it to the given filename.
// It returns an error if any operation fails.
// If the image download is successful, the image is saved to the file and the file is closed.
// If there is any failure in downloading the image, creating the file, or copying the image data to the file,
// an error is returned and all resources are cleaned up.
func DownloadImage(url, filename string) error {
	// Get the image data from the URL.
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download image from URL %s: %w", url, err)
	}
	// Ensure the response body is closed after the function returns.
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("Failed to close response body:", err)
		}
	}()

	// Create the image file.
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	// Ensure the file is closed after the function returns.
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Failed to close file:", err)
		}
	}()

	// Copy the image data to the file.
	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("failed to copy image data to file %s: %w", filename, err)
	}

	return nil
}
