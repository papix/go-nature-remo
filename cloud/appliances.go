package natureremocloud

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Note. I want to use signals only
// please add if you need...
type Appliance struct {
	Signals []Signal `json:"signals"`
}

type Signal struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (c *Client) GetAppliances() ([]*Appliance, error) {
	fmt.Println(c.urlFor("/appliances").String())
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
