package natureremocloud

import (
	"encoding/json"
	"net/http"
)

type Appliance struct {
	ID              string          `json:"id"`
	Nickname        string          `json:"nickname"`
	Image           string          `json:"image"`
	Type            string          `json:"type"`
	Device          Device          `json:"device"`
	Model           Model           `json:"model"`
	CurrentSettings CurrentSettings `json:"settings"`
	Aircon          Aircon          `json:"aircon"`
	Signals         []Signal        `json:"signals"`
}

type Model struct {
	ID           string `json:"id"`
	Manufacturer string `json:"manufacturer"`
	RemoteName   string `json:"remote_name"`
	Name         string `json:"name"`
	Image        string `json:"image"`
}

func (c *Client) GetAppliances() ([]*Appliance, error) {
	req, err := http.NewRequest("GET", c.urlFor("/appliances").String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var appliances []*Appliance
	err = json.NewDecoder(resp.Body).Decode(&appliances)

	if err != nil {
		return nil, err
	}

	return appliances, nil
}
