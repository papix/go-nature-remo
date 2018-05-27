package natureremolocal

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Signal struct {
	Freq   int32   `json:"freq"`
	Data   []int32 `json:"data"`
	Format string  `json:"format"`
}

func NewSignal(freq int32, data []int32, format string) *Signal {
	return &Signal{
		Freq:   freq,
		Data:   data,
		Format: format,
	}
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
	if err := json.NewDecoder(resp.Body).Decode(&signal); err != nil {
		return nil, err
	}

	return signal, nil
}

func (c *Client) PostMessage(signal *Signal) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(signal); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.urlFor("/messages").String(), buf)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return err
	}

	return nil
}
