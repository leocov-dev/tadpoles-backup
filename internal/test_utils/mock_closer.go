package test_utils

import "io"

type MockCloser struct {
	io.Reader
}

func (MockCloser) Close() error { return nil }

func NewMockCloser(r io.Reader) io.ReadCloser {
	return MockCloser{r}
}
