package datadeliverystatus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	BaseURL string
}

type Client struct {
	Config *Config
	HTTP   *http.Client
}

type DataDeliveryStatus interface {
	patch([]byte, string) (string, error)
	Update(string, string) (string, error)
	Error(string, string, string) (string, error)
}

func (client *Client) patch(payload []byte, url string) (string, error) {
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return "", err
	}

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// These functions don't seem to deal with whitespaces and empty strings, not sure if this is intentional
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
