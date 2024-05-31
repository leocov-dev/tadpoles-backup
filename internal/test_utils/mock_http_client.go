package test_utils

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type MockClient struct {
	body io.ReadCloser
}

func NewMockClient(val string) *MockClient {
	return &MockClient{
		body: NewMockCloser(strings.NewReader(val)),
	}
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       m.body,
	}, nil
}

func (m *MockClient) Get(url string) (*http.Response, error) { return m.Do(nil) }
func (m *MockClient) Post(url string, bodyType string, body io.Reader) (*http.Response, error) {
	return m.Do(nil)
}
func (m *MockClient) PostForm(url string, data url.Values) (*http.Response, error) { return m.Do(nil) }
