package directory

import (
	"os/user"
	"strings"

	"github.com/spf13/afero"
)

var fs afero.Fs = afero.NewOsFs()

func SetFs(fileSystem afero.Fs) {
	fs = fileSystem
}

func ExpandHomeDir(dir string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return strings.Replace(dir, "~", usr.HomeDir, 1), nil
}

func CountSeasonDirectories(dir string) (int, error) {
	dir, err := ExpandHomeDir(dir)
	if err != nil {
		return 0, err
	}

	files, err := afero.ReadDir(fs, dir)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, file := range files {
		fileName := file.Name()
		if file.IsDir() && strings.HasPrefix(fileName, "Season") {
			count++
		}
	}
	return count, nil
}
