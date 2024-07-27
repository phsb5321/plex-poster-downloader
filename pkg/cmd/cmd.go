package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"

	"plex-poster-downloader/pkg/config"
	"plex-poster-downloader/pkg/directory"
	"plex-poster-downloader/pkg/poster"
	"plex-poster-downloader/pkg/unsplash"
)

var rootCmd = &cobra.Command{
	Use:   "plex-poster-downloader <directory>",
	Short: "Plex Poster Downloader generates posters for TV show seasons",
	RunE:  execute,
	Args:  cobra.ExactArgs(1),
}

func Execute() error {
	return rootCmd.Execute()
}

func execute(_ *cobra.Command, args []string) error {
	config.Init()

	dir, err := directory.ExpandHomeDir(args[0])
	if err != nil {
		return fmt.Errorf("expanding home directory: %w", err)
	}

	numSeasons, err := directory.CountSeasonDirectories(dir)
	if err != nil {
		return fmt.Errorf("counting season directories: %w", err)
	}

	if numSeasons == 0 {
		logrus.Warn("No season directories found. Generating only the main poster.")
		numSeasons = 1
	}

	unsplashClient := unsplash.NewClient(config.GetUnsplashAccessKey())

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

	if err := poster.GenerateSeasonPosters(dir, numSeasons, unsplashClient, bar.Increment); err != nil {
		return fmt.Errorf("generating season posters: %w", err)
	}

	p.Wait()
	logrus.Info("Season posters generated successfully")
	return nil
}
