package natureremocloud

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var responseBody = `[
  {
    "id": "abcd1234-ef56-gh78-ij90-klmnop123456",
    "device": {
      "id": "qrst7890-uv12-xy34-za56-bcdefg789012",
      "name": "example",
      "temperature_offset": 0,
      "humidity_offset": 0,
      "created_at": "2018-09-25T13:10:02.147Z",
      "updated_at": "2018-09-25T13:10:02.147Z",
      "firmware_version": "Remo/1.0.62-gabbf5bd"
    },
    "model": null,
    "nickname": "light",
    "image": "ico_light",
    "type": "IR",
    "settings": null,
    "aircon": null,
    "signals": [
      {
        "id": "4321abcd-09ef-87gh-65ij-432109kjlmno",
        "name": "switch",
        "image": "ico_on"
      }
    ]
  }
]
`

var applianceHandler = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Authorization") != "Bearer dummy-token" {
		fmt.Fprintln(res, "Authorization header should contains passed API token")
	}

	if req.URL.Path != "/1/appliances" {
		fmt.Fprintf(res, "Request path '%s' is invalid", req.URL.Path)
	}

	res.Header().Set("Content-Type", "application/json")
	io.WriteString(res, responseBody)
})

func TestGetAppliance(t *testing.T) {
	ts := httptest.NewServer(applianceHandler)
	client := NewClientWithOption("dummy-token", ts.URL)
	defer ts.Close()

	var expects = []*Appliance{
		&Appliance{
			ID: "abcd1234-ef56-gh78-ij90-klmnop123456",
			Device: Device{
				ID:                "qrst7890-uv12-xy34-za56-bcdefg789012",
				Name:              "example",
				TemperatureOffset: 0,
				HumidityOffset:    0,
				CreatedAt:         "2018-09-25T13:10:02.147Z",
				UpdatedAt:         "2018-09-25T13:10:02.147Z",
				FirmwareVersion:   "Remo/1.0.62-gabbf5bd",
			},
			Nickname: "light",
			Image:    "ico_light",
			Type:     "IR",
			Signals: []Signal{
				Signal{
					ID:    "4321abcd-09ef-87gh-65ij-432109kjlmno",
					Name:  "switch",
					Image: "ico_on",
				},
			},
		},
	}

	res, err := client.GetAppliances()
	if err != nil {
		t.Errorf("GetAppliance() has ended with error '%s'", err)
	}

	for i, expect := range expects {
		if diff := cmp.Diff(res[i], expect); diff != "" {
			t.Errorf("GetAppliance result ['%d'] differs:\n%s", i, diff)
		}
	}
}
