package utils

import (
	"fmt"
	"io"
	"net/http"
)

type RequestError struct {
	Response *http.Response
	Message  string
}

func (e *RequestError) Error() string {
	defer CloseWithLog(e.Response.Body)
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

func NewRequestError(resp *http.Response, message string) *RequestError {
	return &RequestError{
		Response: resp,
		Message:  message,
	}
}

func PrintErrorList(errorMsgs []string) {
	if errorMsgs != nil {
		WriteError("Errors", "")
		for i, e := range errorMsgs {
			WriteErrorSub.Write(fmt.Sprint(i+1), e)
		}
		fmt.Println("")
	}
}
