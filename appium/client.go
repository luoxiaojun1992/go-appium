package appium

import (
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
