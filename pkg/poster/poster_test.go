package poster_test

import (
	"plex-poster-downloader/pkg/unsplash"
)

// MockUnsplashClient is a mock implementation of the Unsplash client
type MockUnsplashClient struct {
	GetRandomPhotoFunc func() (*unsplash.Photo, error)
	TrackDownloadFunc  func(*unsplash.Photo) error
}

func (m *MockUnsplashClient) GetRandomPhoto() (*unsplash.Photo, error) {
	return m.GetRandomPhotoFunc()
}

func (m *MockUnsplashClient) TrackDownload(photo *unsplash.Photo) error {
	return m.TrackDownloadFunc(photo)
}

// NewMockUnsplashClient creates a new mock Unsplash client
func NewMockUnsplashClient() *MockUnsplashClient {
	return &MockUnsplashClient{
		GetRandomPhotoFunc: func() (*unsplash.Photo, error) {
			return &unsplash.Photo{
				URLs: unsplash.URLs{
					Regular: "http://example.com/image.png",
				},
				Links: unsplash.Links{
					DownloadLocation: "http://example.com/download",
				},
			}, nil
		},
		TrackDownloadFunc: func(*unsplash.Photo) error {
			return nil
		},
	}
}
