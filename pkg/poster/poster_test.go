// File: poster_test.go
// Directory: pkg/poster

// Package poster_test contains unit tests for the poster package.
package poster_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"plex-poster-downloader/pkg/poster"
)

// Constants for testing.
const (
	numSeasons = 3
	baseUrl    = "http://example.com/image.png"
	dir        = "./testdata" // Directory to save the downloaded images.
)

// TestGenerateSeasonPosters verifies that posters can be generated for a given number of seasons and base URL.
func TestGenerateSeasonPosters(t *testing.T) {
	t.Run("Given the number of seasons, a base URL, and a directory", func(t *testing.T) {
		// Ensure the testdata directory exists.
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		err := poster.GenerateSeasonPosters(dir, numSeasons, baseUrl, func() {
			// The progress callback function is an empty function to avoid unnecessary output in the test.
		})

		t.Run("When GenerateSeasonPosters is called", func(t *testing.T) {
			t.Run("Then there should be no error", func(t *testing.T) {
				assert.NoError(t, err, "Error generating season posters")
			})

			t.Run("And a general poster should be generated", func(t *testing.T) {
				filename := fmt.Sprintf("%s/poster.png", dir)
				// Remove the generated file after the test.
				defer os.Remove(filename)

				// Check that the file exists.
				_, err := os.Stat(filename)
				assert.False(t, os.IsNotExist(err), "General poster was not generated")
			})

			t.Run("And each season should have a generated poster", func(t *testing.T) {
				for i := 1; i <= numSeasons; i++ {
					filename := fmt.Sprintf("%s/season%02d-poster.png", dir, i)
					// Remove the generated file after the test.
					defer os.Remove(filename)

					// Check that the file exists.
					_, err := os.Stat(filename)
					assert.False(t, os.IsNotExist(err), "Poster for season %d was not generated", i)
				}
			})
		})
	})
}

// TestGenerateSeasonPostersError verifies that an error is returned when an invalid URL is provided.
func TestGenerateSeasonPostersError(t *testing.T) {
	invalidBaseUrl := "https://invalid-url"
	err := poster.GenerateSeasonPosters(dir, numSeasons, invalidBaseUrl, func() {
		// The progress callback function is an empty function to avoid unnecessary output in the test.
	})
	assert.Error(t, err, "Expected error when generating posters with an invalid URL")
}
