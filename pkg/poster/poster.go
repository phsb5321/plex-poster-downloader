// File: poster/poster.go
package poster

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"plex-poster-downloader/pkg/downloader"
	"plex-poster-downloader/pkg/imageprovider"
)

func GenerateSeasonPosters(dir string, numSeasons int, client *imageprovider.Client, progressCb func()) error {
	var wg sync.WaitGroup
	errors := make(chan error, numSeasons+1)
	rateLimiter := time.NewTicker(time.Second)
	defer rateLimiter.Stop()

	// Get unique photos for all seasons at once
	photos, err := client.GetUniqueRandomPhotos(numSeasons + 1)
	if err != nil {
		return fmt.Errorf("error getting unique random photos: %w", err)
	}

	for i := 0; i <= numSeasons; i++ {
		wg.Add(1)
		go func(season int) {
			defer wg.Done()

			<-rateLimiter.C // Wait for rate limiter

			var posterFilename string
			if season == 0 {
				posterFilename = filepath.Join(dir, "poster.png")
			} else {
				posterFilename = filepath.Join(dir, fmt.Sprintf("season%02d-poster.png", season))
			}

			// Check if the poster already exists
			if _, err := os.Stat(posterFilename); err == nil {
				progressCb()
				return // Skip if poster already exists
			}

			photo := photos[season]

			if err := downloader.DownloadImage(photo.Src.Large, posterFilename); err != nil {
				errors <- fmt.Errorf("error downloading poster for season %d: %w", season, err)
				return
			}

			if err := client.TrackDownload(photo); err != nil {
				errors <- fmt.Errorf("error tracking download for season %d: %w", season, err)
				return
			}

			progressCb()
		}(i)
	}

	go func() {
		wg.Wait()
		close(errors)
	}()

	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}
