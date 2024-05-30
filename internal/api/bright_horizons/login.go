package bright_horizons

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"tadpoles-backup/internal/http_utils"
	"tadpoles-backup/internal/utils"
	"time"
)

func fetchRequestVerificationToken(client *http.Client, rvtUrl *url.URL) (string, error) {

	resp, err := client.Get(rvtUrl.String())
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	r, _ := regexp.Compile("<input.*name=\"__RequestVerificationToken\".*value=\"([a-zA-Z0-9-_]+)\".*/>")

	return r.FindString(string(body)), nil
}

func login(
	client *http.Client,
	loginUrl *url.URL,
	username, password, rvt string,
) (string, error) {
	resp, err := client.PostForm(
		loginUrl.String(),
		url.Values{
			"username":                   {username},
			"password":                   {password},
			"__RequestVerificationToken": {rvt},
			"benefitid":                  {"5"},
			"fstargetid":                 {"1"},
			"usemrType":                  {"0"},
		},
	)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	return resp.Header.Get("Location"), nil
}

func startSaml(client *http.Client, loginUrl *url.URL, redirectVal string) (string, string, error) {
	resp, err := client.Get(fmt.Sprintf("%s%s", loginUrl.String(), redirectVal))
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", err
	}

	actionRx, _ := regexp.Compile("<form.*action=\"([a-zA-Z0-9-_]+)\".*/")
	samlRx, _ := regexp.Compile("<input.*name=\"SAMLResponse\".*value=\"([a-zA-Z0-9-_]+)\".*/>")

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	page := string(body)

	return actionRx.FindString(page), samlRx.FindString(page), nil
}

func finishSaml(client *http.Client, action, samlResponse string) (string, error) {
	resp, err := client.PostForm(
		action,
		url.Values{
			"SAMLResponse": {samlResponse},
		},
	)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "acs" {
			return cookie.Value, nil
		}
	}
	return "", errors.New("BrightHorizons Token acs Cookie not found")
}

type tokenResponse struct {
	Token string `json:"token"`
}

func exchangeToken(client *http.Client, tokenUrl *url.URL, bhToken string) (string, error) {
	req, err := http.NewRequest("GET", tokenUrl.String(), nil)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"Authorization": {fmt.Sprintf("Bearer %s", bhToken)},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	var data tokenResponse
	err = json.Unmarshal(body, &data)

	return data.Token, err
}

func admit(client *http.Client, admitUrl *url.URL, token string, cookieFile string) (*time.Time, error) {
	redirectUrl := admitUrl
	redirectUrl.RawQuery = url.Values{
		"jwt": {token},
	}.Encode()

	resp, err := client.Get(redirectUrl.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return http_utils.SerializeResponseCookies(cookieFile, resp)
}
