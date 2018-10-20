package natureremocloud

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var signalResponseBody = map[string]string{
	"light": `[
  {
    "id": "4321abcd-09ef-87gh-65ij-432109kjlmno",
    "name": "ON",
    "image": "ico_on"
  },
  {
    "id": "0123pqrs-45tu-67vw-89xy-012345zabcde",
    "name": "OFF",
    "image": "ico_off"
  }
]
`,
	"aircon": "[]",
	"tv": `[
  {
    "id": "0001abcd-01ef-23gh-45ij-678901kjlmno",
    "name": "1",
    "image": "ico_1"
  },
  {
    "id": "0002abcd-01ef-23gh-45ij-678901kjlmno",
    "name": "2",
    "image": "ico_2"
  },
  {
    "id": "0003abcd-01ef-23gh-45ij-678901kjlmno",
    "name": "3",
    "image": "ico_3"
  },
  {
    "id": "0004abcd-01ef-23gh-45ij-678901kjlmno",
    "name": "4",
    "image": "ico_4"
  },
  {
    "id": "0005abcd-01ef-23gh-45ij-678901kjlmno",
    "name": "5",
    "image": "ico_5"
  },
  {
    "id": "0006abcd-01ef-23gh-45ij-678901kjlmno",
    "name": "6",
    "image": "ico_6"
  },
  {
    "id": "0007abcd-01ef-23gh-45ij-678901kjlmno",
    "name": "7",
    "image": "ico_7"
  },
  {
    "id": "0008abcd-01ef-23gh-45ij-678901kjlmno",
    "name": "8",
    "image": "ico_8"
  },
  {
    "id": "0009abcd-01ef-23gh-45ij-678901kjlmno",
    "name": "9",
    "image": "ico_9"
  },
  {
    "id": "0010abcd-01ef-23gh-45ij-678901kjlmno",
    "name": "10",
    "image": "ico_10"
  },
  {
    "id": "0011abcd-01ef-23gh-45ij-678901kjlmno",
    "name": "11",
    "image": "ico_11"
  },
  {
    "id": "0012abcd-01ef-23gh-45ij-678901kjlmno",
    "name": "12",
    "image": "ico_12"
  }
]
`,
}

var signalHandler = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Authorization") != "Bearer dummy-token" {
		fmt.Fprintln(res, "Authorization header should contains passed API token")
	}

	if !strings.Contains(req.URL.Path, "/signals") {
		fmt.Fprintf(res, "Request path '%s' is invalid", req.URL.Path)
	}

	applianceID := strings.Split(req.URL.Path, "/")[3]
	res.Header().Set("Content-Type", "application/json")
	io.WriteString(res, signalResponseBody[applianceID])
})

var sendSignalHandler = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Authorization") != "Bearer dummy-token" {
		fmt.Fprintln(res, "Authorization header should contains passed API token")
	}

	splitPath := strings.Split(req.URL.Path, "/")
	if splitPath[2] != "signals" || splitPath[4] != "send" {
		fmt.Fprintf(res, "Request path '%s' is invalid", req.URL.Path)
	}

	res.Header().Set("Content-Type", "application/json")
	re := regexp.MustCompile("signal[0-3]")
	errorResponse := `{
  "code": 404001,
  "message": "Not Found"
}`
	if re.MatchString(splitPath[3]) {
		io.WriteString(res, "{}")
	} else {
		http.Error(res, errorResponse, http.StatusNotFound)
	}

})

func TestGetSignals(t *testing.T) {
	ts := httptest.NewServer(signalHandler)
	client := NewClientWithOption("dummy-token", ts.URL)
	defer ts.Close()

	var tests = []struct {
		applianceID string
		expect      []*Signal
	}{
		{"light", []*Signal{
			&Signal{
				ID:    "4321abcd-09ef-87gh-65ij-432109kjlmno",
				Name:  "ON",
				Image: "ico_on",
			},
			&Signal{
				ID:    "0123pqrs-45tu-67vw-89xy-012345zabcde",
				Name:  "OFF",
				Image: "ico_off",
			},
		}},
		{"aircon", []*Signal{}},
		{"tv", []*Signal{
			&Signal{
				ID:    "0001abcd-01ef-23gh-45ij-678901kjlmno",
				Name:  "1",
				Image: "ico_1",
			},
			&Signal{
				ID:    "0002abcd-01ef-23gh-45ij-678901kjlmno",
				Name:  "2",
				Image: "ico_2",
			},
			&Signal{
				ID:    "0003abcd-01ef-23gh-45ij-678901kjlmno",
				Name:  "3",
				Image: "ico_3",
			},
			&Signal{
				ID:    "0004abcd-01ef-23gh-45ij-678901kjlmno",
				Name:  "4",
				Image: "ico_4",
			},
			&Signal{
				ID:    "0005abcd-01ef-23gh-45ij-678901kjlmno",
				Name:  "5",
				Image: "ico_5",
			},
			&Signal{
				ID:    "0006abcd-01ef-23gh-45ij-678901kjlmno",
				Name:  "6",
				Image: "ico_6",
			},
			&Signal{
				ID:    "0007abcd-01ef-23gh-45ij-678901kjlmno",
				Name:  "7",
				Image: "ico_7",
			},
			&Signal{
				ID:    "0008abcd-01ef-23gh-45ij-678901kjlmno",
				Name:  "8",
				Image: "ico_8",
			},
			&Signal{
				ID:    "0009abcd-01ef-23gh-45ij-678901kjlmno",
				Name:  "9",
				Image: "ico_9",
			},
			&Signal{
				ID:    "0010abcd-01ef-23gh-45ij-678901kjlmno",
				Name:  "10",
				Image: "ico_10",
			},
			&Signal{
				ID:    "0011abcd-01ef-23gh-45ij-678901kjlmno",
				Name:  "11",
				Image: "ico_11",
			},
			&Signal{
				ID:    "0012abcd-01ef-23gh-45ij-678901kjlmno",
				Name:  "12",
				Image: "ico_12",
			},
		}},
	}

	for _, test := range tests {
		res, err := client.GetSignals(test.applianceID)
		if err != nil {
			t.Errorf("GetSignals() has ended with error '%s'", err)
		}

		if diff := cmp.Diff(res, test.expect); diff != "" {
			t.Errorf("GetSignals result ['%s'] differs:\n%s", res, diff)
		}
	}
}

func TestSendSignal(t *testing.T) {
	ts := httptest.NewServer(signalHandler)
	client := NewClientWithOption("dummy-token", ts.URL)
	defer ts.Close()

	var tests = []struct {
		signalID string
		expect   bool
	}{
		{"signal1", true},
		{"signal2", true},
		{"signal3", true},
		{"signal4", false},
		{"dummy-signal", false},
	}

	for _, test := range tests {
		_, err := client.SendSignal(test.signalID)
		if err != nil && test.expect {
			t.Errorf("SendSignal() has ended with error '%s'", err)
		}
	}
}
