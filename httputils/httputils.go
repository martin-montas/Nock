package httputils

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"

type HTTPClient struct {
	client    *http.Client
	userAgent string
}

func NewHTTPClient() *HTTPClient {
	client := &http.Client{
		Timeout: 10 * time.Second, // total request timeout
		Transport: &http.Transport{
			MaxIdleConns:        200,              // total idle connections
			MaxIdleConnsPerHost: 100,              // idle per-host (boosts reuse)
			MaxConnsPerHost:     100,              // hard cap of connections per host
			IdleConnTimeout:     90 * time.Second, // how long to keep idle conns alive
			DisableKeepAlives:   false,            // keep connection open (reused)
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	return &HTTPClient{
		client:    client,
		userAgent: userAgent,
	}
}

func (c *HTTPClient) Get(url string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("Request creation failed:", err)
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println("HTTP request failed:", err)
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}
