package auroradnsclient

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const defaultBaseURL = "https://api.auroradns.eu"

const (
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
)

// Option Type of a client option
type Option func(*Client) error

// Client is a client for accessing the Aurora DNS API
type Client struct {
	baseURL    *url.URL
	UserAgent  string
	httpClient *http.Client
}

// NewClient instantiates a new client
func NewClient(httpClient *http.Client, opts ...Option) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	client := &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}

	for _, opt := range opts {
		err := opt(client)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (c *Client) newRequest(method, resource string, body io.Reader) (*http.Request, error) {
	u, err := c.baseURL.Parse(resource)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set(contentTypeHeader, contentTypeJSON)

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	err = checkResponse(resp)
	if err != nil {
		return resp, err
	}

	if v == nil {
		return resp, nil
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, fmt.Errorf("failed to read body: %v", err)
	}

	err = json.Unmarshal(raw, v)
	if err != nil {
		return resp, fmt.Errorf("unmarshaling error: %v: %s", err, string(raw))
	}

	return resp, nil
}

func checkResponse(resp *http.Response) error {
	if resp.StatusCode < 400 {
		return nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err == nil && data != nil {
		errorResponse := new(ErrorResponse)
		err = json.Unmarshal(data, errorResponse)
		if err != nil {
			return fmt.Errorf("unmarshaling ErrorResponse error: %v: %s", err.Error(), string(data))
		}

		return errorResponse
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

// WithBaseURL Allows to define a custom base URL
func WithBaseURL(rawBaseURL string) func(*Client) error {
	return func(client *Client) error {
		if len(rawBaseURL) == 0 {
			return nil
		}

		baseURL, err := url.Parse(rawBaseURL)
		if err != nil {
			return err
		}

		client.baseURL = baseURL
		return nil
	}
}
