package natureremocloud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
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
}

func TestUrlFor(t *testing.T) {
	client := NewClientWithOption("dummy-token", "https://example.com/ignored/path")
	expect := fmt.Sprintf("https://example.com/%d/correct/endpoint", defaultAPIVersion)

	if url := client.urlFor("/correct/endpoint").String(); url != expect {
		t.Errorf("urlFor should be '%s' but '%s'", expect, url)
	}
}
