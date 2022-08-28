package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var (
	mutex  sync.Mutex
	client *http.Client
)

type Method string

const (
	MethodGet    Method = "GET"
	MethodPost          = "POST"
	MethodPut           = "PUT"
	MethodDelete        = "DELETE"
)

func SetClient(c *http.Client) {
	mutex.Lock()
	client = c
	mutex.Unlock()
}

func Get(uri string, res any, headers ...http.Header) error {
	req, err := http.NewRequest(string(MethodGet), uri, nil)
	if err != nil {
		return err
	}

	// set request header
	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header[key] = value
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteArray, res)
	if err != nil {
		return fmt.Errorf("error: %s, url=%s, json: [%s]", err, uri, string(byteArray))
	}

	return nil
}

func Post(uri string, data any, res any, headers ...http.Header) (int, error) {
	return Execute(uri, MethodPost, data, res, headers...)
}

func Put(uri string, data any, res any, headers ...http.Header) (int, error) {
	return Execute(uri, MethodPut, data, res, headers...)
}

func Delete(uri string, data any, res any, headers ...http.Header) (int, error) {
	return Execute(uri, MethodDelete, data, res, headers...)
}

func Execute(uri string, method Method, data any, res any, headers ...http.Header) (int, error) {
	// create json
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest(string(method), uri, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}

	// set request header
	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header[key] = value
		}
	} else {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(byteArray, res)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("error: %s, url=%s, json: [%s]", err, uri, string(byteArray))
	}

	return resp.StatusCode, nil
}
