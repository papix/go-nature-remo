package natureremocloud

import (
	"encoding/json"
	"net/http"
)

type Device struct {
	ID                string       `json:"id"`
	Name              string       `json:"name"`
	TemperatureOffset int32        `json:"temperature_offset"`
	HumidityOffset    int32        `json:"humidity_offset"`
	CreatedAt         string       `json:"created_at"`
	UpdatedAt         string       `json:"updated_at"`
	FirmwareVersion   string       `json:"firmware_version"`
	NewestEvents      NewestEvents `json:"newest_events"`
}

type NewestEvents struct {
	Temperature Temperature `json:"te"`
	Humidity    Humidity    `json:"hu"`
}

type Temperature struct {
	Value     float64 `json:"val"`
	CreatedAt string  `json:"created_at"`
}

type Humidity struct {
	Value     int32  `json:"val"`
	CreatedAt string `json:"created_at"`
}

func (c *Client) GetDevices() ([]*Device, error) {
	req, err := http.NewRequest("GET", c.urlFor("/devices").String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var devices []*Device
	err = json.NewDecoder(resp.Body).Decode(&devices)

	if err != nil {
		return nil, err
	}

	return devices, nil
}
