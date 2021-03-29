package datadeliverystatus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Config struct {
	BaseURL string
}

type Client struct {
	Config *Config
	HTTP   *http.Client
}

func (client *Client) patch(payload []byte, url string) (string, error) {
	req, err := client.HTTP.NewRequest(http.MethodPatch, url, bytes.NewBuffe(payload))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return "", err
	}

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (client *Client) Update(state, fileName string) (string, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"state": state,
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/v1/state/%s", client.Config.BaseURL, fileName)

	return client.patch(payload, url)
}

func (client *Client) Error(state, fileName, errorMessage string) (string, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"state":      state,
		"error_info": errorMessage,
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/v1/state/%s", client.Config.BaseURL, fileName)

	return client.patch(payload, url)
}
