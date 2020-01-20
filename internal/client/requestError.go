package client

import (
	"fmt"
	"net/http"
)

type RequestError struct {
	Response *http.Response
	Message  string
}

func (e *RequestError) Error() string {
	return e.Message
}

func NewRequestError(resp *http.Response) *RequestError {
	return &RequestError{
		Response: resp,
		Message:  fmt.Sprintf("[Error] %s: %s", resp.Request.Method, resp.Request.URL.Path),
	}
}
