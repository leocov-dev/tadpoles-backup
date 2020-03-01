package client

import (
	"encoding/json"
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
)

var (
	homeDir, _         = os.UserHomeDir()
	TadpolesCookieFile = path.Join(homeDir, ".tadpole-backup-cookie")

	Jar, _    = cookiejar.New(nil)
	ApiClient = &http.Client{Jar: Jar}

	tadpolesHost   = "https://www.tadpoles.com"
	TadpolesUrl, _ = url.Parse(tadpolesHost)

	LoginEndpoint = fmt.Sprintf("%s://%s/auth/login", TadpolesUrl.Scheme, TadpolesUrl.Host)

	apiV1               = fmt.Sprintf("%s://%s/remote/v1", TadpolesUrl.Scheme, TadpolesUrl.Host)
	EventsEndpoint      = fmt.Sprintf("%s/events", apiV1)
	AttachmentsEndpoint = fmt.Sprintf("%s/obj_attachment", apiV1)
	AdmitEndpoint       = fmt.Sprintf("%s/athome/admit", apiV1)
	ParametersEndpoint  = fmt.Sprintf("%s/parameters", apiV1)
)

// load cookies from serialized json on disk if able.
func DeserializeCookies() {
	var storedCookies []*http.Cookie
	if utils.FileExists(TadpolesCookieFile) {

		// Open our jsonFile
		jsonFile, _ := os.Open(TadpolesCookieFile)
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		err := json.Unmarshal(byteValue, &storedCookies)

		if err != nil {
			log.Debug("Failed to deserialize cookies...", err)
			return
		}
		log.Debug("Deserialized cookies from file")
	}
	// load cookies to cookie jar that api client will use
	Jar.SetCookies(TadpolesUrl, storedCookies)
}

func SerializeCookies() {
	cookiesData := Jar.Cookies(TadpolesUrl)
	jsonString, _ := json.MarshalIndent(cookiesData, "", "  ")

	err := ioutil.WriteFile(TadpolesCookieFile, jsonString, 0644)

	if err != nil {
		log.Debug("Failed to serialize cookies to file...", err)
		return
	}

	log.Debug("Serialize cookies successful")
}
