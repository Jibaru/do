package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jibaru/do/internal/types"
)

type HttpClient interface {
	Do(doFile types.DoFile) (*http.Response, error)
}

type httpClient struct {
	client *http.Client
}

func NewHttpClient(client *http.Client) HttpClient {
	return &httpClient{
		client,
	}
}

func (h *httpClient) Do(doFile types.DoFile) (*http.Response, error) {
	// Replace params
	url := doFile.Do.URL
	for key, value := range doFile.Do.Params {
		placeholder := fmt.Sprintf(":%s", key)
		url = strings.Replace(url, placeholder, fmt.Sprintf("%v", value), -1)
	}

	req, err := http.NewRequest(doFile.Do.Method, url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range doFile.Do.Headers {
		req.Header.Add(key, fmt.Sprintf("%v", value))
	}

	query := req.URL.Query()
	for key, value := range doFile.Do.Query {
		query.Add(key, fmt.Sprintf("%v", value))
	}
	req.URL.RawQuery = query.Encode()

	if doFile.Do.Body != "" {
		bodyBytes, err := json.Marshal(doFile.Do.Body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
	}

	return h.client.Do(req)
}
