package downloader_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"plex-poster-downloader/pkg/downloader"
)

func TestDownloadImage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		imageData, _ := ioutil.ReadFile("plex-poster-downloader/test/testdata/test_image.png")
		w.Write(imageData)
	}))
	defer ts.Close()

	tempFile := "test/temp_downloaded_image.png"
	defer os.Remove(tempFile)

	err := downloader.DownloadImage(ts.URL, tempFile)
	require.NoError(t, err, "Error downloading image")

	_, err = os.Stat(tempFile)
	assert.False(t, os.IsNotExist(err), "Image file was not downloaded")
}
