package client

import (
	"encoding/json"
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/config"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	log "github.com/sirupsen/logrus"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

var (
	jar, _    = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.CookieJarList})
	ApiClient = newApiClient()

	tadpolesHost   = "https://www.tadpoles.com"
	TadpolesUrl, _ = url.Parse(tadpolesHost)

	LoginEndpoint = fmt.Sprintf("%s://%s/auth/login", TadpolesUrl.Scheme, TadpolesUrl.Host)

	apiV1               = fmt.Sprintf("%s://%s/remote/v1", TadpolesUrl.Scheme, TadpolesUrl.Host)
	EventsEndpoint      = fmt.Sprintf("%s/events", apiV1)
	AttachmentsEndpoint = fmt.Sprintf("%s/obj_attachment", apiV1)
	AdmitEndpoint       = fmt.Sprintf("%s/athome/admit", apiV1)
	ParametersEndpoint  = fmt.Sprintf("%s/parameters", apiV1)
)

func newApiClient() *http.Client {
	deserializeCookies()
	return &http.Client{Jar: jar}
}

// load cookies from serialized json on disk if able.
func deserializeCookies() {
	var storedCookies []*http.Cookie
	if utils.FileExists(config.TadpolesCookieFile) {

		// Open our jsonFile
		jsonFile, _ := os.Open(config.TadpolesCookieFile)
		defer utils.CloseWithLog(jsonFile)

		byteValue, _ := ioutil.ReadAll(jsonFile)
		err := json.Unmarshal(byteValue, &storedCookies)

		if err != nil {
			log.Debug("Failed to deserialize cookies...", err)
			return
		}
		log.Debug(fmt.Sprintf("Deserialized cookies from file: %s", config.TadpolesCookieFile))
	}
	// load cookies to cookie jar that api client will use
	jar.SetCookies(TadpolesUrl, storedCookies)
}

func SerializeResponseCookies(response *http.Response) {
	cookiesData := response.Cookies()
	jsonString, _ := json.MarshalIndent(cookiesData, "", "  ")

	err := ioutil.WriteFile(config.TadpolesCookieFile, jsonString, 0600)

	if err != nil {
		log.Debug("Failed to serialize cookies to file...", err)
		return
	}

	log.Debug("Serialize cookies successful")
}
