package poster

import (
	"fmt"
	"sync"

	"plex-poster-downloader/pkg/downloader"
)

func GenerateSeasonPosters(dir string, numSeasons int, baseURL string, progressCb func()) error {
	var wg sync.WaitGroup
	errors := make(chan error)

	for i := 0; i <= numSeasons; i++ {
		wg.Add(1)
		go func(season int) {
			defer wg.Done()

			// Define the poster filename.
			var posterFilename string
			if season == 0 {
				// For the general poster image.
				posterFilename = fmt.Sprintf("%s/poster.png", dir)
			} else {
				// For season-specific posters.
				posterFilename = fmt.Sprintf("%s/season%02d-poster.png", dir, season)
			}

			// Download the poster from the baseURL.
			if err := downloader.DownloadImage(baseURL, posterFilename); err != nil {
				errors <- fmt.Errorf("error downloading poster for season %d: %w", season, err)
				return
			}

			// Call the progress callback function after each successful download.
			progressCb()
		}(i)
	}

	// Wait for all downloads to finish.
	go func() {
		wg.Wait()
		close(errors)
	}()

	// Return the first error that occurs, if any.
	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}
