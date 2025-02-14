// src/services/httpclient.go
package httpclient

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// HTTPClient is a simple reusable client wrapper
type HTTPClient struct {
	client  *http.Client
	BaseURL string
}

// NewHTTPClient initializes an HTTPClient with a custom timeout, base URL, etc.
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		BaseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second, // used so that if the server is taking too much time, the request can be terminated from the client side
			//here the server is the endpint being hit and client is our server.
		},
	}
}

func (hc *HTTPClient) Get(path string) ([]byte, error) {
	fullURL := fmt.Sprintf("%s%s", hc.BaseURL, path)

	resp, err := hc.client.Get(fullURL)
	if err != nil {
		log.Printf("Error in HTTPClient.Get: %v", err)
		return nil, fmt.Errorf("GET request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Printf("Error in HTTPClient.Get: %v, GET request returned status code %v ", err, resp.StatusCode)
		return nil, fmt.Errorf("GET request returned status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error in HTTPClient.Get: %v, failed to read response body ", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return bodyBytes, nil
}
