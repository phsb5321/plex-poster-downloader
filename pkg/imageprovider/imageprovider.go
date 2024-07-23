package imageprovider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const pexelsAPIURL = "https://api.pexels.com/v1"

type Client struct {
	AccessKey     string
	HTTPClient    *http.Client
	RateLimiter   *time.Ticker
	RequestsCount int
	mu            sync.Mutex
}

type Photo struct {
	ID           int    `json:"id"`
	Src          Src    `json:"src"`
	URL          string `json:"url"`
	Photographer string `json:"photographer"`
}

type Src struct {
	Original string `json:"original"`
	Large    string `json:"large"`
	Medium   string `json:"medium"`
	Small    string `json:"small"`
}

func NewClient(accessKey string) *Client {
	return &Client{
		AccessKey:     accessKey,
		HTTPClient:    &http.Client{},
		RateLimiter:   time.NewTicker(time.Second / 5), // Limit to 5 requests per second
		RequestsCount: 0,
	}
}

func (c *Client) GetRandomPhoto() (*Photo, error) {
	c.mu.Lock()
	<-c.RateLimiter.C
	c.RequestsCount++
	count := c.RequestsCount
	c.mu.Unlock()

	logrus.Debugf("Making request #%d", count)

	endpoint := fmt.Sprintf("%s/search", pexelsAPIURL)

	params := url.Values{}
	params.Add("query", "abstract colorful meaningful")
	params.Add("orientation", "landscape")
	params.Add("per_page", "1")

	req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.AccessKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		logrus.Warn("Rate limit hit, waiting for 10 seconds before retrying")
		time.Sleep(time.Second * 10)
		return c.GetRandomPhoto()
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Photos []Photo `json:"photos"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Photos) > 0 {
		return &result.Photos[0], nil
	}

	return nil, fmt.Errorf("no photos found")
}

func (c *Client) TrackDownload(photo *Photo) error {
	if photo == nil || photo.Src.Original == "" {
		return fmt.Errorf("invalid photo or download location")
	}

	// We don't need to actually download the image for Pexels, just return nil
	return nil
}

func (c *Client) GetUniqueRandomPhotos(count int) ([]*Photo, error) {
	var (
		photos    = make([]*Photo, 0, count)
		seenIDs   = make(map[int]bool)
		mu        sync.Mutex
		wg        sync.WaitGroup
		errChan   = make(chan error, count)
		photoChan = make(chan *Photo, count)
	)

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			for {
				photo, err := c.GetRandomPhoto()
				if err != nil {
					logrus.Errorf("Error fetching photo for index %d: %v", index, err)
					errChan <- err
					return
				}

				mu.Lock()
				if !seenIDs[photo.ID] {
					seenIDs[photo.ID] = true
					photoChan <- photo
					logrus.Debugf("Added photo %d for index %d", photo.ID, index)
					mu.Unlock()
					return
				}
				mu.Unlock()
				logrus.Debugf("Duplicate photo %d for index %d, retrying", photo.ID, index)
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(photoChan)
		close(errChan)
	}()

	for photo := range photoChan {
		photos = append(photos, photo)
	}

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	if len(photos) < count {
		return nil, fmt.Errorf("could only fetch %d unique photos out of %d requested", len(photos), count)
	}

	return photos, nil
}
