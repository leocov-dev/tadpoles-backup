package tadpoles_api

import (
	"github.com/leocov-dev/tadpoles-backup/config"
	"net/http"
	"net/http/cookiejar"
)

var (
	Jar, _    = cookiejar.New(nil)
	apiClient = &http.Client{Jar: Jar}
)

func init() {
	var storedCookies []*http.Cookie
	Jar.SetCookies(config.TadpolesUrl, storedCookies)
}
