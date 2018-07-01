package natureremocloud

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultBaseURL    = "https://api.nature.global/"
	defaultAPIVersion = 1
	defaultUserAgent  = "go-nature-remo"
	apiRequestTimeout = 30 * time.Second

	limit      int64 = 30
	hLimit           = "X-Rate-Limit-Limit"
	hRemaining       = "X-Rate-Limit-Remaining"
	hReset           = "X-Rate-Limit-Reset"
)

type Limit struct {
	Limit     int64
	Remaining int64
	Reset     int64
}

type Client struct {
	BaseURL     *url.URL
	APIVersion  int32
	AccessToken string
	UserAgent   string
	HTTPClient  *http.Client
	Limit       *Limit
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
	c.Limit = &Limit{Limit: limit}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	if v, err := strconv.ParseInt(resp.Header.Get(hLimit), 10, 64); err == nil {
		c.Limit.Limit = v
	}

	if v, err := strconv.ParseInt(resp.Header.Get(hRemaining), 10, 64); err == nil {
		c.Limit.Remaining = v
	}
	if v, err := strconv.ParseInt(resp.Header.Get(hReset), 10, 64); err == nil {
		c.Limit.Reset = v
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		err := fmt.Errorf("API result failed: %s. It will recover in %d secs",
			resp.Status, c.Limit.Reset-time.Now().Unix())
		return resp, err
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
