package http

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"
	"varq/config"
)

var Cfg *config.Config
var ErrHTTPNotOk = errors.New("HTTP response with status code not 200 OK")

func Get(url string) ([]byte, error) {
	client := http.Client{
		Timeout: time.Duration(Cfg.HTTPClient.Timeout) * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", Cfg.HTTPClient.UserAgent)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, ErrHTTPNotOk
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
