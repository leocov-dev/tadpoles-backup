package test_utils

import "net/url"

func MockUrl(val string) *url.URL {
	u, _ := url.Parse(val)
	return u
}
