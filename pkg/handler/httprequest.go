package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

func HttpRequset(method, url string, jsonBody any) error {
	client := http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	req.Header = http.Header{
		"Authorization": {"Basic hhgSdnV6aW5rYXMx"},
		"Content-Type":  {"application/json; charset=utf-8"},
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}
	bodyString := string(body)
	var inventory any
	if err := json.Unmarshal([]byte(bodyString), &inventory); err != nil {
		return err
	}
	jsonStr, err := json.Marshal(inventory)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(jsonStr), &jsonBody)
	if err != nil {
		return err
	}
	return nil
}
