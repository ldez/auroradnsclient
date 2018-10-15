package auroradnsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Record a DNS record
type Record struct {
	ID         string `json:"id,omitempty"`
	RecordType string `json:"type"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	TTL        int    `json:"ttl,omitempty"`
}

// GetRecords returns a list of all records in given zone
func (c *Client) GetRecords(zoneID string) ([]Record, *http.Response, error) {
	relativeURL := fmt.Sprintf("/zones/%s/records", zoneID)

	req, err := c.newRequest(http.MethodGet, relativeURL, nil)
	if err != nil {
		return nil, nil, err
	}

	var respData []Record
	resp, err := c.do(req, &respData)
	if err != nil {
		return nil, resp, err
	}

	return respData, resp, nil
}

// CreateRecord creates a new record in given zone
func (c *Client) CreateRecord(zoneID string, record Record) (*Record, *http.Response, error) {
	body, err := json.Marshal(record)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshall request body: %v", err)
	}

	relativeURL := fmt.Sprintf("/zones/%s/records", zoneID)

	req, err := c.newRequest(http.MethodPost, relativeURL, bytes.NewReader(body))
	if err != nil {
		return nil, nil, err
	}

	respData := new(Record)
	resp, err := c.do(req, respData)
	if err != nil {
		return nil, resp, err
	}

	return respData, resp, nil
}

// RemoveRecord removes a record corresponding to a particular id in a given zone
func (c *Client) RemoveRecord(zoneID string, recordID string) (bool, *http.Response, error) {
	relativeURL := fmt.Sprintf("/zones/%s/records/%s", zoneID, recordID)

	req, err := c.newRequest(http.MethodDelete, relativeURL, nil)
	if err != nil {
		return false, nil, err
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return false, resp, err
	}

	return true, resp, nil
}
