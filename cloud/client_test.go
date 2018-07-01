package natureremocloud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set(hLimit, fmt.Sprintf("%d", limit))
		res.Header().Set(hRemaining, fmt.Sprintf("%d", (limit-1)))
		res.Header().Set(hReset, fmt.Sprintf("%d", time.Now().Unix()+300))

		if req.Header.Get("Authorization") != "Bearer dummy-token" {
			t.Errorf("Authorization header should contains passed API token")
		}

		if h := req.Header.Get("User-Agent"); h != defaultUserAgent {
			t.Errorf("User-Agent shoud be '%s' but %s", defaultUserAgent, h)
		}
	}))
	defer ts.Close()

	client := NewClientWithOption("dummy-token", ts.URL)

	req, _ := http.NewRequest("GET", client.urlFor("/").String(), nil)
	client.Request(req)

	if client.Limit.Limit != limit {
		t.Errorf("Limit.Limit should be %d, but %d",
			limit, client.Limit.Limit)
	}
	if client.Limit.Remaining > client.Limit.Limit {
		t.Errorf("Limit.Remaining should be 0 <= x < %d, but %d",
			client.Limit.Limit, client.Limit.Remaining)
	}

	if n := time.Now().Unix(); client.Limit.Reset < n {
		t.Errorf("Limit.Reset should be in the future, but %v",
			client.Limit.Reset)
	}
}

func TestTooManyRequests(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set(hLimit, fmt.Sprintf("%d", limit))
		res.Header().Set(hRemaining, "0")
		res.Header().Set(hReset, fmt.Sprintf("%d", time.Now().Unix()+300))
		res.WriteHeader(http.StatusTooManyRequests)

	}))
	defer ts.Close()

	client := NewClientWithOption("dummy-token", ts.URL)

	req, _ := http.NewRequest("GET", client.urlFor("/").String(), nil)
	resp, err := client.Request(req)
	if err == nil {
		t.Error("Client should have an error")
	}

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Client shoud have StatusCode %d, but %d",
			http.StatusTooManyRequests, resp.StatusCode)
	}

	if client.Limit.Remaining != 0 {
		t.Errorf("Client should have no remaining, but %d",
			client.Limit.Remaining)
	}
}

func TestUrlFor(t *testing.T) {
	client := NewClientWithOption("dummy-token", "https://example.com/ignored/path")
	expect := fmt.Sprintf("https://example.com/%d/correct/endpoint", defaultAPIVersion)

	if url := client.urlFor("/correct/endpoint").String(); url != expect {
		t.Errorf("urlFor should be '%s' but '%s'", expect, url)
	}
}
