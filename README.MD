# Plex Poster Downloader

Plex Poster Downloader is a tool to generate and download season-specific posters for TV shows.

## Features

- Generate season posters based on the number of seasons.
- Download posters from a given base URL.
- Progress callback functionality to track the poster generation process.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine.

### Prerequisites

You need to have [Go](https://golang.org/) installed on your machine. This project was built using Go version 1.16, but it should work with other versions as well.

### Installation

To install the project, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/username/plex-poster-downloader.git
   ```

2. Navigate to the project directory:

   ```bash
   cd plex-poster-downloader
   ```

3. Build the project:

   ```bash
   go build
   ```

### Usage

To use the Plex Poster Downloader, follow these steps:

1. Set up your configuration file (config.yml) in the config directory with your base URL.

2. Run the following command in your terminal:

   ```bash
   ./plex-poster-downloader
   ```

## Running the tests

To run the tests, use the following command:

```bash
go test ./...
```

## Built With

- [Go](https://golang.org/) - The programming language used.
- [Cobra](https://github.com/spf13/cobra) - A library for creating powerful modern CLI applications.
- [Viper](https://github.com/spf13/viper) - A complete configuration solution for Go applications.
- [mpb](https://github.com/vbauerster/mpb) - A progress bar library for Go.

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
