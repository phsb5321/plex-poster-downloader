package poster

import (
	"fmt"
	"sync"
	"time"

	"plex-poster-downloader/pkg/downloader"
	"plex-poster-downloader/pkg/unsplash"
)

func GenerateSeasonPosters(dir string, numSeasons int, unsplashClient *unsplash.Client, progressCb func()) error {
	var wg sync.WaitGroup
	errors := make(chan error)
	rateLimiter := time.NewTicker(time.Second)
	defer rateLimiter.Stop()

	for i := 0; i <= numSeasons; i++ {
		wg.Add(1)
		go func(season int) {
			defer wg.Done()

			<-rateLimiter.C // Wait for rate limiter

			photo, err := unsplashClient.GetRandomPhoto()
			if err != nil {
				errors <- fmt.Errorf("error getting random photo for season %d: %w", season, err)
				return
			}

			var posterFilename string
			if season == 0 {
				posterFilename = fmt.Sprintf("%s/poster.png", dir)
			} else {
				posterFilename = fmt.Sprintf("%s/season%02d-poster.png", dir, season)
			}

			if err := downloader.DownloadImage(photo.URLs.Regular, posterFilename); err != nil {
				errors <- fmt.Errorf("error downloading poster for season %d: %w", season, err)
				return
			}

			if err := unsplashClient.TrackDownload(photo); err != nil {
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
