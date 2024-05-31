package bright_horizons

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"net/http"
	"net/url"
	"strings"
	"tadpoles-backup/internal/interfaces"
	"tadpoles-backup/internal/utils"
)

func checkLogin(client interfaces.HttpClient, checkUrl *url.URL) error {
	_, err := fetchDependents(client, checkUrl)
	return err
}

// fetchRequestVerificationToken
// parse the CSRF token out of the initial login HTML page so we can make a
// non-interactive login call
func fetchRequestVerificationToken(client interfaces.HttpClient, rvtUrl *url.URL) (string, error) {

	resp, err := client.Get(rvtUrl.String())
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", utils.NewRequestError(resp, "error fetching initial bh page")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	doc, err := htmlquery.Parse(strings.NewReader(string(body)))

	node := htmlquery.FindOne(doc, "//input[@name='__RequestVerificationToken']")
	found := htmlquery.SelectAttr(node, "value")

	if found == "" {
		return "", errors.New("error finding RVT in page")
	}

	return found, nil
}

// login
// initial login form for bright horizons users, requires a request verification
// token used to counter CSRF (which we are doing...)
func login(
	client interfaces.HttpClient,
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
func startSaml(client interfaces.HttpClient, samlUrl *url.URL) (string, string, error) {
	resp, err := client.Get(samlUrl.String())
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", utils.NewRequestError(resp, "error with SAML 1st step")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	doc, err := htmlquery.Parse(strings.NewReader(string(body)))

	actionNode := htmlquery.FindOne(doc, "//form")
	action := htmlquery.SelectAttr(actionNode, "action")
	samlNode := htmlquery.FindOne(doc, "//input[@name='SAMLResponse']")
	saml := htmlquery.SelectAttr(samlNode, "value")

	if action == "" || saml == "" {
		return "", "", errors.New("error finding SAML response")
	}

	return action, saml, nil
}

// finishSaml
// second step in SAML process should yield the bright horizons API token
// in a response cookie
func finishSaml(client interfaces.HttpClient, action, samlResponse string) (string, error) {
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
func exchangeToken(client interfaces.HttpClient, tokenUrl *url.URL, bhToken string) (string, error) {
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

func admitRedirect(client interfaces.HttpClient, redirectUrl *url.URL) ([]*http.Cookie, error) {
	resp, err := client.Get(redirectUrl.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, utils.NewRequestError(resp, "error validating jwt tadpoles token")
	}

	return resp.Cookies(), nil
}
