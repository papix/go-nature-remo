package natureremocloud

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL    = "https://api.nature.global/"
	defaultAPIVersion = 1
	defaultUserAgent  = "go-nature-remo"
	apiRequestTimeout = 30 * time.Second
)

type Client struct {
	BaseURL     *url.URL
	APIVersion  int32
	AccessToken string
	UserAgent   string
	HTTPClient  *http.Client
}

func NewClient(token string) *Client {
	u, _ := url.Parse(defaultBaseURL)
	return &Client{
		BaseURL:     u,
		APIVersion:  defaultAPIVersion,
		AccessToken: token,
		UserAgent:   defaultUserAgent,
		HTTPClient:  &http.Client{},
	}
}

func NewClientWithOption(token string, baseURL string) *Client {
	u, _ := url.Parse(baseURL)
	return &Client{
		BaseURL:     u,
		APIVersion:  defaultAPIVersion,
		AccessToken: token,
		UserAgent:   defaultUserAgent,
		HTTPClient:  &http.Client{},
	}
}

func (c *Client) Request(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
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

	newURL.Path = fmt.Sprintf("/%d%s", c.APIVersion, endpoint)

	return newURL
}

func closeResponse(resp *http.Response) {
	if resp != nil {
		resp.Body.Close()
	}
}
