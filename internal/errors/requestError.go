package errors

import (
	"net/http"
)

type RequestError struct {
	Response *http.Response
	Message  string
}

func (e *RequestError) Error() string {
	return e.Message
}
