package poster_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"plex-poster-downloader/pkg/poster"
)

const (
	numSeasons = 3
	baseUrl    = "http://example.com/image.png"
)

func TestGenerateSeasonPosters(t *testing.T) {
	t.Run("Given the number of seasons and a base URL", func(t *testing.T) {
		err := poster.GenerateSeasonPosters(numSeasons, baseUrl, func() {
			// Empty function to avoid unnecessary output in the test
		})

		t.Run("When GenerateSeasonPosters is called", func(t *testing.T) {
			t.Run("There should be no error", func(t *testing.T) {
				assert.NoError(t, err, "Error generating season posters")
			})

			t.Run("Each season should have a generated poster", func(t *testing.T) {
				for i := 1; i <= numSeasons; i++ {
					filename := fmt.Sprintf("season%02d-poster.png", i)
					defer os.Remove(filename)

					_, err := os.Stat(filename)
					assert.False(t, os.IsNotExist(err), "Poster for season %d was not generated", i)
				}
			})
		})
	})
}
