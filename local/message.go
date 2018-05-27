package natureremolocal

import (
	"encoding/json"
	"net/http"
)

type Signal struct {
	Freq   int32   `json:"freq"`
	Data   []int32 `json:"data"`
	Format string  `json:"format"`
}

func (c *Client) GetMessage() (*Signal, error) {
	req, err := http.NewRequest("GET", c.urlFor("/messages").String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var signal *Signal
	err = json.NewDecoder(resp.Body).Decode(&signal)

	if err != nil {
		return nil, err
	}

	return signal, nil
}
