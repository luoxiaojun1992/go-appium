package appium

type Element struct {
    ID        string
    Sess   *Session
}

func (e *Element) Click() error {
    url := fmt.Sprintf("%s/element/%s/click", e.Sess.URL, e.ID)
    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        return err
    }

    res, err := e.Sess.Client.Do(req)
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

func (e *Element) SendKeys(text string) error {
    url := fmt.Sprintf("%s/element/%s/value", e.Sess.URL, e.ID)
    params := map[string]interface{}{
        "value": strings.Split(text, ""),
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

    res, err := e.Sess.Client.Do(req)
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
