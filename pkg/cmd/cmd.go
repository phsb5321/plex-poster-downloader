// File: cmd/cmd.go
package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"

	"plex-poster-downloader/pkg/config"
	"plex-poster-downloader/pkg/directory"
	"plex-poster-downloader/pkg/imageprovider"
	"plex-poster-downloader/pkg/poster"
)

var rootCmd = &cobra.Command{
	Use:   "plex-poster-downloader",
	Short: "Plex Poster Downloader generates posters for TV show seasons",
	Run:   execute,
	Args:  cobra.ExactArgs(1),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execute(cmd *cobra.Command, args []string) {
	config.Init()

	dir := args[0]

	numSeasons, err := directory.CountSeasonDirectories(dir)
	if err != nil {
		logrus.Fatalf("Error counting season directories: %v", err)
	}

	if numSeasons == 0 {
		logrus.Warn("No season directories found. Generating only the main poster.")
		numSeasons = 1 // Set to 1 to generate at least the main poster
	}

	accessKey := config.GetImageProviderAccessKey()
	client := imageprovider.NewClient(accessKey)

	p := mpb.New(mpb.WithWidth(60))
	bar := p.AddBar(int64(numSeasons),
		mpb.PrependDecorators(
			decor.Name("Generating season posters: "),
			decor.CountersNoUnit("%d / %d"),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WCSyncSpace),
		),
	)

	dir, err = directory.ExpandHomeDir(dir)
	if err != nil {
		logrus.Fatalf("Error expanding home directory: %v", err)
	}

	err = poster.GenerateSeasonPosters(dir, numSeasons, client, func() { bar.Increment() })
	if err != nil {
		logrus.Fatalf("Error generating season posters: %v", err)
	}

	p.Wait()
	logrus.Info("Season posters generated successfully")
}
