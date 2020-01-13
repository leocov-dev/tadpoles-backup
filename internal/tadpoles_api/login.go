package tadpoles_api

import (
	"fmt"
	"net/http"
	"net/url"
)

func PostLogin(email string, password string) {
	resp, err := http.PostForm("http://example.com/form",
		url.Values{"email": {email}, "password": {password}, "service": {"tadpoles"}},
	)
	fmt.Printf("Response: %+v\n", resp)
	fmt.Printf("Error   : %+v\n", err)
}

func storeAuthCookie(cookie string) {

}
