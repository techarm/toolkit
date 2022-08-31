package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Method string

const (
	MethodGet    Method = "GET"
	MethodPost   Method = "POST"
	MethodPut    Method = "PUT"
	MethodDelete Method = "DELETE"
)

type Client struct {
	Header     http.Header
	HttpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		HttpClient: http.DefaultClient,
	}
}

func (c *Client) GetString(uri string, headers ...http.Header) (string, int, error) {
	req, err := http.NewRequest(string(MethodGet), uri, nil)
	if err != nil {
		return "", 0, err
	}

	// set request header
	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header[key] = value
		}
	} else if len(c.Header) > 0 {
		for key, value := range c.Header {
			req.Header[key] = value
		}
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", 0, err
	}

	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	return string(byteArray), resp.StatusCode, nil
}

func (c *Client) Get(uri string, res any, headers ...http.Header) (int, error) {
	req, err := http.NewRequest(string(MethodGet), uri, nil)
	if err != nil {
		return 0, err
	}

	// set request header
	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header[key] = value
		}
	} else if len(c.Header) > 0 {
		for key, value := range c.Header {
			req.Header[key] = value
		}
	}

	resp, err := c.HttpClient.Do(req)
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

func (c *Client) Post(uri string, data any, res any, headers ...http.Header) (int, error) {
	return c.Execute(uri, MethodPost, data, res, headers...)
}

func (c *Client) Put(uri string, data any, res any, headers ...http.Header) (int, error) {
	return c.Execute(uri, MethodPut, data, res, headers...)
}

func (c *Client) Delete(uri string, data any, res any, headers ...http.Header) (int, error) {
	return c.Execute(uri, MethodDelete, data, res, headers...)
}

func (c *Client) Execute(uri string, method Method, data any, res any, headers ...http.Header) (int, error) {
	var body io.Reader

	// if data is set, convert it to json
	if data != nil {
		// create json
		jsonData, err := json.Marshal(data)
		if err != nil {
			return 0, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(string(method), uri, body)
	if err != nil {
		return 0, err
	}

	// set request header
	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header[key] = value
		}
	} else if len(c.Header) > 0 {
		for key, value := range c.Header {
			req.Header[key] = value
		}
	} else {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if res != nil {
		err = json.Unmarshal(byteArray, res)
		if err != nil {
			return resp.StatusCode, fmt.Errorf("error: %s, url=%s, json: [%s]", err, uri, string(byteArray))
		}
	}

	return resp.StatusCode, nil
}
