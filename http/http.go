package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"respdb/config"
	"time"
)

var Cfg *config.Config

func Get(url string) ([]byte, error) {
	timeout := 120
	userAgent := "Test"
	if Cfg != nil {
		timeout = Cfg.HTTPClient.Timeout
		userAgent = Cfg.HTTPClient.UserAgent
	}

	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP status code %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
