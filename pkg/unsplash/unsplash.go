package unsplash

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const unsplashAPIURL = "https://api.unsplash.com"

type Client struct {
	AccessKey  string
	HTTPClient *http.Client
}

type Photo struct {
	ID    string `json:"id"`
	URLs  URLs   `json:"urls"`
	Links Links  `json:"links"`
}

type URLs struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
}

type Links struct {
	Self             string `json:"self"`
	HTML             string `json:"html"`
	Download         string `json:"download"`
	DownloadLocation string `json:"download_location"`
}

func NewClient(accessKey string) *Client {
	return &Client{
		AccessKey:  accessKey,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) GetRandomPhoto() (*Photo, error) {
	endpoint := fmt.Sprintf("%s/photos/random", unsplashAPIURL)

	params := url.Values{}
	params.Add("query", "abstract colorful meaningful")
	params.Add("orientation", "landscape")
	params.Add("content_filter", "high")

	req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", c.AccessKey))
	req.Header.Set("Accept-Version", "v1")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var photo Photo
	if err := json.NewDecoder(resp.Body).Decode(&photo); err != nil {
		return nil, err
	}

	return &photo, nil
}

func (c *Client) TrackDownload(photo *Photo) error {
	if photo == nil || photo.Links.DownloadLocation == "" {
		return fmt.Errorf("invalid photo or download location")
	}

	req, err := http.NewRequest("GET", photo.Links.DownloadLocation, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", c.AccessKey))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to track download: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
