package request

import (
	"net/http"
	"sync"
)

var (
	root  *Client
	mutex sync.Mutex
)

func SetClient(c *http.Client) {
	mutex.Lock()
	root.HttpClient = c
	mutex.Unlock()
}

func GetString(uri string, headers ...http.Header) (string, int, error) {
	return root.GetString(uri, headers...)
}

func Get(uri string, res any, headers ...http.Header) (int, error) {
	return root.Get(uri, res, headers...)
}

func Post(uri string, data any, res any, headers ...http.Header) (int, error) {
	return root.Post(uri, data, res, headers...)
}

func Put(uri string, data any, res any, headers ...http.Header) (int, error) {
	return root.Put(uri, data, res, headers...)
}

func Delete(uri string, data any, res any, headers ...http.Header) (int, error) {
	return root.Delete(uri, data, res, headers...)
}
