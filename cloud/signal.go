package natureremocloud

import (
	"encoding/json"
	"net/http"
)

type Signal struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (c *Client) GetSignals(applianceID string) ([]*Signal, error) {
	req, err := http.NewRequest("GET", c.urlFor("/appliances/"+applianceID+"/signals").String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var signals []*Signal
	err = json.NewDecoder(resp.Body).Decode(&signals)

	if err != nil {
		return nil, err
	}

	return signals, nil
}
