package appium

type Session struct {
    ID            string
    PlatformName  string
    PlatformVer   string
    DeviceName    string
    App           string
    Automation    string
    WebDriverAddr string
    Client        *http.Client
}

func (s *Session) Start() error {
    url := fmt.Sprintf("%s/session", s.WebDriverAddr)
    params := map[string]interface{}{
        "capabilities": map[string]interface{}{
            "platformName":  s.PlatformName,
            "platformVer":   s.PlatformVer,
            "deviceName":    s.DeviceName,
            "app":           s.App,
            "automationName": s.Automation,
        },
    }

    data, err := json.Marshal(params)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")

    res, err := s.Client.Do(req)
    if err != nil {
        return err
    }

    defer res.Body.Close()

    var result map[string]interface{}
    err = json.NewDecoder(res.Body).Decode(&result)
    if err != nil {
        return err
    }

    if result["status"].(int) != 0 {
        return fmt.Errorf(result["value"].(string))
    }

    s.ID = result["sessionId"].(string)

    return nil
}

func (s *Session) Stop() error {
    url := fmt.Sprintf("%s/session/%s", s.WebDriverAddr, s.ID)

    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return err
    }

    res, err := s.Client.Do(req)
    if err != nil {
        return err
    }

    defer res.Body.Close()

    var result map[string]interface{}
    err = json.NewDecoder(res.Body).Decode(&result)
    if err != nil {
        return err
    }

    if result["status"].(int) != 0 {
        return fmt.Errorf(result["value"].(string))
    }

    return nil
}

func (s *Session) FindElement(using string, value string) (*Element, error) {
    url := fmt.Sprintf("%s/element", s.URL)
    params := map[string]interface{}{
        "using": using,
        "value": value,
    }

    data, err := json.Marshal(params)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")

    res, err := s.Client.Do(req)
    if err != nil {
        return nil, err
    }

    defer res.Body.Close()

    var result map[string]interface{}
    err = json.NewDecoder(res.Body).Decode(&result)
    if err != nil {
        return nil, err
    }

    if result["status"].(int) != 0 {
        return nil, fmt.Errorf(result["value"].(string))
    }

    elementID := result["value"].(map[string]interface{})["ELEMENT"].(string)
    return &Element{ID: elementID, Session: s}, nil
}

func (s *Session) Status() (map[string]interface{}, error) {
    url := fmt.Sprintf("%s/status", s.URL)

    res, err := s.Client.Get(url)
    if err != nil {
        return nil, err
    }

    defer res.Body.Close()

    var result map[string]interface{}
    err = json.NewDecoder(res.Body).Decode(&result)
    if err != nil {
        return nil, err
    }

    return result, nil
}

func (s *Session) Log(logType string) ([]map[string]interface{}, error) {
    url := fmt.Sprintf("%s/log/%s", s.URL, logType)

    res, err := s.Client.Get(url)
    if err != nil {
        return nil, err
    }

    defer res.Body.Close()

    var result []map[string]interface{}
    err = json.NewDecoder(res.Body).Decode(&result)
    if err != nil {
        return nil, err
    }

    return result, nil
}
