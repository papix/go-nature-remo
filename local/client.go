package natureremolocal

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultUserAgent  = "go-nature-remo"
	apiRequestTimeout = 30 * time.Second
)

type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	u, _ := url.Parse(baseURL)
	return &Client{
		BaseURL:    u,
		UserAgent:  defaultUserAgent,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) Request(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", c.UserAgent)

	client := c.HTTPClient
	client.Timeout = apiRequestTimeout

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, fmt.Errorf("API result failed: %s", resp.Status)
	}
	return resp, nil
}

func (c *Client) urlFor(endpoint string) *url.URL {
	newURL, err := url.Parse(c.BaseURL.String())
	if err != nil {
		panic("invalid url passed")
	}

	newURL.Path = endpoint

	return newURL
}

func closeResponse(resp *http.Response) {
	if resp != nil {
		resp.Body.Close()
	}
}
