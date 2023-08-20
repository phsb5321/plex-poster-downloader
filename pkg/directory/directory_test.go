package directory_test

import (
	"os/user"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"plex-poster-downloader/pkg/directory"
)

func TestExpandHomeDir(t *testing.T) {
	expandedDir, err := directory.ExpandHomeDir("~")
	assert.NoError(t, err)

	usr, _ := user.Current()
	expectedDir := usr.HomeDir

	assert.Equal(t, expectedDir, expandedDir)
}

func TestCountSeasonDirectories(t *testing.T) {
	testFs := afero.NewMemMapFs()
	directory.SetFs(testFs)

	// Create directories in the in-memory filesystem.
	_ = testFs.MkdirAll("/testdir/Season1", 0755)
	_ = testFs.MkdirAll("/testdir/Season2", 0755)
	_ = testFs.MkdirAll("/testdir/NotASeason", 0755)

	// Count the season directories.
	count, err := directory.CountSeasonDirectories("/testdir")
	assert.NoError(t, err)

	assert.Equal(t, 2, count)
}
