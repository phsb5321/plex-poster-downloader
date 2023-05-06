package poster

import (
	"fmt"

	"plex-poster-downloader/pkg/downloader"
)

func GenerateSeasonPosters(numSeasons int, baseUrl string, progressCb func()) error {
	for i := 1; i <= numSeasons; i++ {
		seasonPosterFilename := fmt.Sprintf("season%02d-poster.png", i)
		err := downloader.DownloadImage(baseUrl, seasonPosterFilename)
		if err != nil {
			return fmt.Errorf("Error downloading season %d poster: %w", i, err)
		}
		progressCb()

	}
	return nil
}
