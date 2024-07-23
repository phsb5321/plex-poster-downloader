// File: poster/poster_test.go
package poster_test

import (
	"plex-poster-downloader/pkg/imageprovider"
)

// MockImageProviderClient is a mock implementation of the ImageProvider client
type MockImageProviderClient struct {
	GetRandomPhotoFunc func() (*imageprovider.Photo, error)
	TrackDownloadFunc  func(*imageprovider.Photo) error
}

func (m *MockImageProviderClient) GetRandomPhoto() (*imageprovider.Photo, error) {
	return m.GetRandomPhotoFunc()
}

func (m *MockImageProviderClient) TrackDownload(photo *imageprovider.Photo) error {
	return m.TrackDownloadFunc(photo)
}

// NewMockImageProviderClient creates a new mock ImageProvider client
func NewMockImageProviderClient() *MockImageProviderClient {
	return &MockImageProviderClient{
		GetRandomPhotoFunc: func() (*imageprovider.Photo, error) {
			return &imageprovider.Photo{
				Src: imageprovider.Src{
					Large: "http://example.com/image.png",
				},
				URL: "http://example.com/download",
			}, nil
		},
		TrackDownloadFunc: func(*imageprovider.Photo) error {
			return nil
		},
	}
}
