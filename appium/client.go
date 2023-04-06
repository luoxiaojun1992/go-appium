package appium

import (
    "net/http"
    "net/url"
)

type Client struct {
	ServerUrl *url.URL
}

func NewClient(serverUrl string) (*Client, error) {
	u, err := url.Parse(serverUrl)
	if err != nil {
		return nil, err
	}

	return &Client{
		ServerUrl: u,
	}, nil
}

func (c *Client) NewSession(opts ...SessionOption) *Session {
	options := &SessionOptions{}

	// Apply the options to the options struct
	for _, opt := range opts {
		opt(options)
	}

	// Create a new session object with the options
	return &Session{
		PlatformName: options.PlatformName,
		PlatformVer:  options.PlatformVer,
		DeviceName:   options.DeviceName,
		App:          options.App,
		Automation:   options.Automation,
		Client:       &http.Client{},
	}
}
