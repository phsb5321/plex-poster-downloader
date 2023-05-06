package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"

	"plex-poster-downloader/pkg/poster"
)

var rootCmd = &cobra.Command{
	Use:   "plex-poster-downloader",
	Short: "Plex Poster Downloader generates posters for TV show seasons",
	Run:   execute,
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("Fatal error reading config file: %s", err)
	}
}

func countSeasonDirectories() (int, error) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return 0, err
	}

	count := 0
	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), "Season") {
			count++
		}
	}
	return count, nil
}

func execute(_ *cobra.Command, _ []string) {
	initConfig()

	numSeasons, err := countSeasonDirectories()
	if err != nil {
		logrus.Fatalf("Error counting season directories: %v", err)
	}

	baseUrl := viper.GetString("baseUrl")

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

	err = poster.GenerateSeasonPosters(numSeasons, baseUrl, func() { bar.Increment() })
	if err != nil {
		logrus.Fatalf("Error generating season posters: %v", err)
	}

	p.Wait()
	logrus.Info("Season posters generated successfully")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
