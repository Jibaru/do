package request

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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
		switch doFile.Do.Body.(type) {
		case types.String:
			val := doFile.Do.Body.(types.String)
			req.Body = io.NopCloser(strings.NewReader(string(val)))
		case types.Map:
			// Map equals to multipart/form-data
			var requestBody bytes.Buffer
			writer := multipart.NewWriter(&requestBody)

			for key, value := range doFile.Do.Body.(types.Map) {
				switch value.(type) {
				case types.String:
					val := value.(types.String)
					err = writer.WriteField(key, string(val))
					if err != nil {
						return nil, NewCanNotDoRequestError(err)
					}
				case types.File:
					typeFile := value.(types.File)

					file, err := os.Open(typeFile.Path)
					if err != nil {
						return nil, NewCanNotDoRequestError(err)
					}
					part, err := writer.CreateFormFile(key, file.Name())
					if err != nil {
						return nil, NewCanNotDoRequestError(err)
					}

					_, err = io.Copy(part, file)
					if err != nil {
						return nil, NewCanNotDoRequestError(err)
					}
				}
			}

			err = writer.Close()
			if err != nil {
				return nil, NewCanNotDoRequestError(err)
			}

			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Body = io.NopCloser(&requestBody)
		}
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
