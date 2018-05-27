package natureremolocal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-Requested-With") != "XMLHttpRequest" {
			t.Errorf("X-Requested-With header is required")
		}

		if h := req.Header.Get("User-Agent"); h != defaultUserAgent {
			t.Errorf("User-Agent shoud be '%s' but %s", defaultUserAgent, h)
		}
	}))
	defer ts.Close()

	client := NewClient(ts.URL)

	req, _ := http.NewRequest("GET", client.urlFor("/").String(), nil)
	client.Request(req)
}

func TestUrlFor(t *testing.T) {
	client := NewClient("https://example.com/ignored/path")
	expect := "https://example.com/correct/endpoint"

	if url := client.urlFor("/correct/endpoint").String(); url != expect {
		t.Errorf("urlFor should be '%s' but '%s'", expect, url)
	}
}
