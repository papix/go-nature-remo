package natureremocloud

import (
	"encoding/json"
	"net/http"
)

func (c *Client) SendSignal(signalID string) (*string, error) {
	req, err := http.NewRequest("POST", c.urlFor("/signals/"+signalID+"/send").String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var rep_msg string
	err = json.NewDecoder(resp.Body).Decode(&rep_msg)

	if err != nil {
		return nil, err
	}

	return &rep_msg, nil
}
