package appium

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
)

type SessionOptions struct {
	PlatformName  string
	PlatformVer   string
	DeviceName    string
	App           string
	Automation    string
}

type SessionOption func(*SessionOptions)

func WithPlatformName(platformName string) SessionOption {
	return func(o *SessionOptions) {
		o.PlatformName = platformName
	}
}

func WithPlatformVer(platformVer string) SessionOption {
	return func(o *SessionOptions) {
		o.PlatformVer = platformVer
	}
}

func WithDeviceName(deviceName string) SessionOption {
	return func(o *SessionOptions) {
		o.DeviceName = deviceName
	}
}

func WithApp(app string) SessionOption {
	return func(o *SessionOptions) {
		o.App = app
	}
}

func WithAutomation(automation string) SessionOption {
	return func(o *SessionOptions) {
		o.Automation = automation
	}
}

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
    url := fmt.Sprintf("%s/element", s.WebDriverAddr)
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
    url := fmt.Sprintf("%s/status", s.WebDriverAddr)

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
    url := fmt.Sprintf("%s/log/%s", s.WebDriverAddr, logType)

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

func (s *Session) Lock(duration int) error {
    url := fmt.Sprintf("%s/appium/device/lock", s.WebDriverAddr)

    data := map[string]int{"seconds": duration}
    body, err := json.Marshal(data)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")

    res, err := s.Client.Do(req)
    if err != nil {
        return err
    }

    defer res.Body.Close()

    return nil
}

func (s *Session) Unlock() error {
    url := fmt.Sprintf("%s/appium/device/unlock", s.WebDriverAddr)

    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        return err
    }

    res, err := s.Client.Do(req)
    if err != nil {
        return err
    }

    defer res.Body.Close()

    return nil
}

func (s *Session) InstallApp(appPath string) error {
    url := fmt.Sprintf("%s/session/%s/appium/device/install_app", s.WebDriverAddr, s.ID)

    file, err := os.Open(appPath)
    if err != nil {
        return err
    }
    
    defer file.Close()

    req, err := http.NewRequest("POST", url, file)
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

func (s *Session) GetPerformanceData(packageName string, dataType string, dataReadTimeout time.Duration) (map[string]interface{}, error) {
    // 构造请求体
    requestData := map[string]interface{}{
        "packageName":    packageName,
        "dataType":       dataType,
        "dataReadTimeout": int(dataReadTimeout / time.Millisecond),
    }
    requestBodyBytes, err := json.Marshal(requestData)
    if err != nil {
        return nil, err
    }

    // 发送请求
    req, err := http.NewRequest("POST", s.WebDriverAddr+"/session/"+s.ID+"/appium/performanceData", bytes.NewReader(requestBodyBytes))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    resp, err := s.Client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // 解析响应体
    var responseData map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&responseData)
    if err != nil {
        return nil, err
    }
    return responseData, nil
}
