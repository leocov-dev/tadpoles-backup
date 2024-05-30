package bright_horizons

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"tadpoles-backup/internal/utils"
)

func checkLogin(client *http.Client, checkUrl *url.URL) error {
	_, err := fetchDependents(client, checkUrl)
	return err
}

// fetchRequestVerificationToken
// parse the CSRF token out of the initial login HTML page so we can make a
// non-interactive login call
func fetchRequestVerificationToken(client *http.Client, rvtUrl *url.URL) (string, error) {

	resp, err := client.Get(rvtUrl.String())
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", utils.NewRequestError(resp, "error fetching initial bh page")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	page := string(body)

	r, _ := regexp.Compile("<input.*name=\"__RequestVerificationToken\".*value=\"([a-zA-Z0-9-_]+)\".*/>")

	return r.FindString(page), nil
}

// login
// initial login form for bright horizons users, requires a request verification
// token used to counter CSRF (which we are doing...)
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
		return "", utils.NewRequestError(resp, "bh login error")
	}

	return resp.Header.Get("Location"), nil
}

// startSaml
// first step in SAML process should yield an "action" url and a SAML Response
// value
func startSaml(client *http.Client, samlUrl *url.URL) (string, string, error) {
	resp, err := client.Get(samlUrl.String())
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", utils.NewRequestError(resp, "error with SAML 1st step")
	}

	actionRx, _ := regexp.Compile("<form.*action=\"([a-zA-Z0-9-_]+)\".*/")
	samlRx, _ := regexp.Compile("<input.*name=\"SAMLResponse\".*value=\"([a-zA-Z0-9-_]+)\".*/>")

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	page := string(body)

	return actionRx.FindString(page), samlRx.FindString(page), nil
}

// finishSaml
// second step in SAML process should yield the bright horizons API token
// in a response cookie
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
		return "", utils.NewRequestError(resp, "error with SAML 2nd step")
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

// exchangeToken
// exchange the bright horizons token for a tadpoles token
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
		return "", utils.NewRequestError(resp, "error exchanging bright horizons token for tadpoles token")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	var data tokenResponse
	err = json.Unmarshal(body, &data)

	return data.Token, err
}

func admitRedirect(client *http.Client, redirectUrl *url.URL) ([]*http.Cookie, error) {
	resp, err := client.Get(redirectUrl.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, utils.NewRequestError(resp, "error validating jwt tadpoles token")
	}

	return resp.Cookies(), nil
}
