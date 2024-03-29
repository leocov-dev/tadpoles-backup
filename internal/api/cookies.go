package api

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"tadpoles-backup/internal/utils"
	"time"
)

// DeserializeCookies load cookies from serialized json on disk if able.
func DeserializeCookies(cookieFile string, hostUrl *url.URL) *cookiejar.Jar {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.CookieJarList})

	if utils.FileExists(cookieFile) {
		// Open our jsonFile
		jsonFile, err := os.Open(cookieFile)
		defer utils.CloseWithLog(jsonFile)
		if err != nil {
			logrus.Debug("Failed to open cookie json file...", err)
			return jar
		}

		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			logrus.Debug("Failed to read cookie json file...", err)
			return jar
		}

		var storedCookies []*http.Cookie
		err = json.Unmarshal(byteValue, &storedCookies)

		if err != nil {
			logrus.Debug("Failed to unmarshal cookies...", err)
			return jar
		}

		logrus.Debug(fmt.Sprintf("Deserialized cookies from file: %s", cookieFile))
		jar.SetCookies(hostUrl, storedCookies)
	}

	// may return empty cookie jar if no serialized file found
	return jar
}

func SerializeResponseCookies(cookieFile string, response *http.Response) (expires *time.Time, err error) {
	cookiesData := response.Cookies()
	jsonString, err := json.MarshalIndent(cookiesData, "", "  ")
	if err != nil {
		logrus.Debug("Failed to marshal cookies...", err)
		return nil, err
	}

	err = os.WriteFile(cookieFile, jsonString, 0600)
	if err != nil {
		logrus.Debug("Failed to write cookies json to file...", err)
		return nil, err
	}

	logrus.Debug("Serialize cookies successful")
	return &cookiesData[0].Expires, nil
}
