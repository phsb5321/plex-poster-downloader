// File: downloader_test.go
// Directory: pkg/downloader

package downloader_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"plex-poster-downloader/pkg/downloader"
)

// TestDownloadImage verifies that the DownloadImage function successfully downloads an image from a URL.
func TestDownloadImage(t *testing.T) {
	// Set up an HTTP test server that serves a mock image
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		imageData, err := ioutil.ReadFile("../../test/testdata/mock_image.png")
		if err != nil {
			t.Fatalf("Error reading test image file: %v", err)
		}
		_, err = w.Write(imageData)
		if err != nil {
			t.Fatalf("Failed to write image data to response: %v", err)
		}
	}))
	defer ts.Close()

	// Create a temporary file to store the downloaded image
	tempFile := "temp_downloaded_image.png"
	defer func() {
		err := os.Remove(tempFile)
		if err != nil {
			t.Logf("Failed to remove temporary file: %s", err)
		}
	}()

	// Download the image using the DownloadImage function
	err := downloader.DownloadImage(ts.URL, tempFile)
	if err != nil {
		t.Fatalf("Failed to download image: %v", err)
	}

	// Check if the image file was downloaded
	_, err = os.Stat(tempFile)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("Image file was not downloaded: %v", err)
		}
		t.Fatalf("Failed to stat downloaded file: %v", err)
	}
}

// TestDownloadImageError verifies that the DownloadImage function returns an error when attempting to download
// an image from an invalid URL.
func TestDownloadImageError(t *testing.T) {
	invalidUrl := "http://invalid-url"
	tempFile := "test/temp_downloaded_image.png"
	defer func() {
		err := os.Remove(tempFile)
		if err != nil {
			t.Logf("Failed to remove temporary file: %s", err)
		}
	}()

	err := downloader.DownloadImage(invalidUrl, tempFile)
	assert.Error(t, err, "Expected error when downloading from invalid URL")
}

// TestDownloadImageCreateFileError verifies that the DownloadImage function returns an error when it fails to create a file.
func TestDownloadImageCreateFileError(t *testing.T) {
	// Set up an HTTP test server that serves a mock image
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		imageData, err := ioutil.ReadFile("../../test/testdata/mock_image.png")
		if err != nil {
			t.Fatalf("Error reading test image file: %v", err)
		}
		_, err = w.Write(imageData)
		if err != nil {
			t.Fatalf("Failed to write image data to response: %v", err)
		}
	}))
	defer ts.Close()

	// Provide an empty filename, which should cause os.Create to return an error
	tempFile := ""
	err := downloader.DownloadImage(ts.URL, tempFile)
	assert.Error(t, err, "Expected error when failing to create a file")
}

// TestDownloadImageCopyError verifies that the DownloadImage function returns an error when it fails to copy the image data
// to the file.
func TestDownloadImageCopyError(t *testing.T) {
	// Set up an HTTP test server that simulates the presence of data by setting the Content-Length header
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Length", "1")
	}))
	defer ts.Close()

	// Create a temporary file to store the downloaded image
	tempFile := "temp_downloaded_image.png"
	defer func() {
		err := os.Remove(tempFile)
		if err != nil {
			t.Logf("Failed to remove temporary file: %s", err)
		}
	}()

	// Temporarily replace http.DefaultClient with a custom client that returns a custom reader
	originalClient := http.DefaultClient
	defer func() { http.DefaultClient = originalClient }()
	http.DefaultClient = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			resp, err := originalClient.Get(ts.URL)
			if err != nil {
				return nil, err
			}

			resp.Body = ioutil.NopCloser(&errorReader{})
			return resp, nil
		}),
	}

	// Download the image using the DownloadImage function
	err := downloader.DownloadImage(ts.URL, tempFile)
	assert.Error(t, err, "Expected error when failing to copy image data to the file")
}

// roundTripperFunc is a helper type that allows implementing the http.RoundTripper interface with a function.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (rtf roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rtf(req)
}

// errorReader is a custom reader that always returns an error when attempting to read data.
type errorReader struct{}

func (er *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}
