package api

import (
	"fmt"
	"io"
	"net/http"
	"tadpoles-backup/internal/utils"
)

type RequestError struct {
	Response *http.Response
	Message  string
}

func (e *RequestError) Error() string {
	defer utils.CloseWithLog(e.Response.Body)
	body, _ := io.ReadAll(e.Response.Body)

	msg := fmt.Sprintf(
		"[Error] %s %s: %s => %s",
		e.Message,
		e.Response.Request.Method,
		e.Response.Request.URL.Path,
		string(body),
	)
	return msg
}

func newRequestError(resp *http.Response, message string) *RequestError {
	return &RequestError{
		Response: resp,
		Message:  message,
	}
}
