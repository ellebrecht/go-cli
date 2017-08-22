package request

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	log "geeny/log"
	url "geeny/net/url"
)

// Builder builds http requests
type Builder struct {
	urlBuilder     *url.Builder
	defaultHeaders map[string]string
}

// NewBuilder returns a new Builder instance
func NewBuilder(baseURL string, defaultHeaders map[string]string) *Builder {
	return &Builder{
		urlBuilder: url.NewBuilder(baseURL),
	}
}

// NewRequest creates a new http request
func (b *Builder) NewRequest(httpMethod string, endpoint string, params map[string]string, headers map[string]string, body []byte) (*http.Request, error) {
	var req *http.Request
	var err error
	if body != nil && len(body) > 0 {
		req, err = http.NewRequest(httpMethod, b.urlBuilder.NewURL(endpoint), bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(httpMethod, b.urlBuilder.NewURL(endpoint), nil)
	}
	if err != nil {
		return nil, err
	}

	if params != nil {
		curParams := req.URL.Query()
		for key, value := range params {
			curParams.Add(key, value)
		}
		req.URL.RawQuery = curParams.Encode()
	}
	for key, value := range b.defaultHeaders {
		req.Header.Set(key, value)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	log.Debugf("Request: %+v", *req)
	return req, nil
}

// NewMultipartRequest creates a new multipart http request
func (b *Builder) NewMultipartRequest(httpMethod string, endpoint string, params map[string]string, headers map[string]string, name string, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(name, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		err = writer.WriteField(key, val)
		if err != nil {
			return nil, err
		}
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(httpMethod, b.urlBuilder.NewURL(endpoint), body)
	for key, value := range b.defaultHeaders {
		req.Header.Set(key, value)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	log.Debugf("Multipart request: %+v", *req)
	return req, err
}
