package auroradnsclient

import (
	"net/http"
)

// Zone a DNS zone
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetZones returns a list of all zones
func (c *Client) GetZones() ([]Zone, *http.Response, error) {
	req, err := c.newRequest(http.MethodGet, "/zones", nil)
	if err != nil {
		return nil, nil, err
	}

	var zones []Zone
	resp, err := c.do(req, &zones)
	if err != nil {
		return nil, resp, err
	}

	return zones, resp, nil
}
