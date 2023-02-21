package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(url, token string) (string, error) {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", url, nil)
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("source %s not found", url))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
