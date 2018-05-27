package natureremolocal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetMessage(t *testing.T) {
	expect := NewSignal(0, []int32{1, 2, 3}, "string")

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/messages" {
			t.Errorf("request URL should be '/messages' but '%s'", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Errorf("request method should be 'GET' but '%s'", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"freq":   expect.Freq,
			"data":   expect.Data,
			"format": expect.Format,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	signal, err := client.GetMessage()

	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	}

	if signal.Freq != expect.Freq {
		t.Errorf("Freq should be %d, but %d", expect.Freq, signal.Freq)
	}

	if signal.Format != expect.Format {
		t.Errorf("Freq should be '%s', but '%s'", expect.Format, signal.Format)
	}

	if reflect.DeepEqual(signal.Data, expect.Data) != true {
		t.Errorf("Wrong data: %v", signal.Data)
	}
}

func TestPostMessage(t *testing.T) {
	expect := NewSignal(0, []int32{1, 2, 3}, "string")

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/messages" {
			t.Errorf("request URL should be '/messages' but '%s'", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Errorf("request method should be 'POST' but '%s'", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)
		signal := &Signal{}
		err := json.Unmarshal(body, signal)
		if err != nil {
			t.Fatalf("request body should be decode as json: %s", string(body))
		}

		if signal.Freq != expect.Freq {
			t.Errorf("Freq should be %d, but %d", expect.Freq, signal.Freq)
		}

		if signal.Format != expect.Format {
			t.Errorf("Freq should be '%s', but '%s'", expect.Format, signal.Format)
		}

		if reflect.DeepEqual(signal.Data, expect.Data) != true {
			t.Errorf("Wrong data: %v", signal.Data)
		}
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	err := client.PostMessage(expect)

	if err != nil {
		t.Errorf("err should be nil, but %v", err)
	}
}
