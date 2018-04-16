package natureremocloud

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type Me struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
}

func (c *Client) GetMe() (*Me, error) {
	req, err := http.NewRequest("GET", c.urlFor("/users/me").String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var me *Me
	err = json.NewDecoder(resp.Body).Decode(&me)

	if err != nil {
		return nil, err
	}

	return me, nil
}

func (c *Client) PostMe(nickname string) (*Me, error) {
	values := url.Values{}
	values.Set("nickname", nickname)

	req, err := http.NewRequest("POST", c.urlFor("/users/me").String(), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var me *Me
	err = json.NewDecoder(resp.Body).Decode(&me)

	if err != nil {
		return nil, err
	}

	return me, nil
}
