package request

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jibaru/do/internal/types"
)

type HttpClient interface {
	Do(doFile types.DoFile) (*types.Response, error)
}

type httpClient struct {
	client *http.Client
}

func NewHttpClient(client *http.Client) HttpClient {
	return &httpClient{
		client,
	}
}

func (h *httpClient) Do(doFile types.DoFile) (*types.Response, error) {
	// Replace params
	url := string(doFile.Do.URL)
	for key, value := range doFile.Do.Params {
		placeholder := fmt.Sprintf(":%s", key)
		beforeReplaceUrl := url
		afterReplaceUrl := strings.Replace(url, placeholder, fmt.Sprintf("%v", value), -1)

		if beforeReplaceUrl == afterReplaceUrl {
			return nil, NewCanNotReplaceParamError(key)
		}

		url = afterReplaceUrl
	}

	req, err := http.NewRequest(string(doFile.Do.Method), url, nil)
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

	if doFile.Do.Body != nil {
		body := string(*doFile.Do.Body)
		req.Body = io.NopCloser(strings.NewReader(body))
	}

	res, err := h.client.Do(req)
	if err != nil {
		return nil, NewCanNotDoRequestError(err)
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, NewCanNotReadResponseBodyError(err)
	}

	// Get response headers
	headers := make(map[string]interface{})
	for key, value := range res.Header {
		headers[key] = value
	}

	return &types.Response{
		StatusCode: res.StatusCode,
		Body:       string(respBody),
		Headers:    headers,
	}, nil
}
