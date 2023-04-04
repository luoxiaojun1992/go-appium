package appium

type AppiumClient struct {
	ServerUrl *url.URL
}

func NewAppiumClient(serverUrl string) (*AppiumClient, error) {
	u, err := url.Parse(serverUrl)
	if err != nil {
		return nil, err
	}

	return &AppiumClient{
		ServerUrl: u,
	}, nil
}
