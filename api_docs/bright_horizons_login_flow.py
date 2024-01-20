"""
Script to test new BrightHorizons login flow.

execute from command line:
$ python bright_horizons_login_flow.py <username> <password>
"""
import argparse
import json
import re
from http import cookiejar
from urllib import error, parse, request

LOGIN_URL = "https://bhlogin.brighthorizons.com"
API_URL = "https://mybrightday.brighthorizons.com"

initial_params = {
    "benefitid":  5,
    "fstargetid": 1,
}

token_pattern = re.compile(
    "<input.*name=\"__RequestVerificationToken\".*value=\"(?P<token>[a-zA-Z0-9-_]+)\".*/>",
    re.MULTILINE,
)

login_denied_pattern = re.compile("We can.*t find that Personal Username and/or Password", re.MULTILINE)
login_locked_pattern = re.compile("An incorrect Personal Username/Password has been entered", re.MULTILINE)


def get_csrf_data():
    # get the request verification token.
    # must call url with benefit/fstargetid params in order to get the correct login form.
    initial_url = f"{LOGIN_URL}?{parse.urlencode(initial_params)}"

    req = request.Request(initial_url)
    with request.urlopen(req) as resp:
        page: str = resp.read().decode('utf-8')
        cj = cookiejar.CookieJar()
        cj.extract_cookies(resp, req)

    # print("\npage ----------\n", page)
    # print(cj)

    # parse the __RequestVerificationToken out of the login form's hidden field
    match = next(token_pattern.finditer(page), None)

    # print(match)

    return match.group("token"), cj


def post_login_form(username: str, password: str, token: str, cj: cookiejar.CookieJar):
    form_data = parse.urlencode(
        {
            # **initial_params,
            "__RequestVerificationToken": token,
            "username":                   username,
            "password":                   password,
            "response":                   "jwt"
        }
    )

    req = request.Request(
        LOGIN_URL,
        data=form_data.encode('utf-8'),
        headers={
            "Content-Type": "application/x-www-form-urlencoded",
        },
        method="POST",
    )
    cj.add_cookie_header(req)

    # print(req.header_items())

    try:
        with request.urlopen(req) as resp:
            page = resp.read().decode('utf-8')
            cj = cookiejar.CookieJar()
            cj.extract_cookies(resp, req)
    except error.HTTPError as e:
        print(e.read().decode('utf-8'))
        raise e

    if re.search(login_denied_pattern, page):
        # print(page)
        raise Exception("incorrect username or password")
    if re.search(login_locked_pattern, page):
        # print(page)
        raise Exception("account is locked")

    print("\npost login page ----------\n", page)
    print("\npost login cookies ----------\n", cj)


def validate_jwt_token(token: str) -> str:
    """ turn a token provided by login flow into an api key """
    api_key = ""

    form_data = parse.urlencode(
        {
            "token": token
        }
    )

    req = request.Request(
        f"{API_URL}/api/v2/auth/jwt/validate",
        data=form_data.encode('utf-8'),
        headers={
            "Content-Type": "application/x-www-form-urlencoded",
        },
    )
    try:
        with request.urlopen(req) as resp:
            body = json.loads(resp.read().decode('utf-8'))
            api_key = body["apiKey"]
    except error.HTTPError as e:
        print(e.read().decode('utf-8'))

    return api_key


def get_user_profile(api_key: str):
    """ get user profile data """
    req = request.Request(
        f"{API_URL}/api/v2/user/profile",
        method="GET",
        headers={
            "Accept":    "application/json",
            "X-Api-Key": api_key
        }
    )
    try:
        with request.urlopen(req, ) as resp:
            print(resp.read())
    except error.HTTPError as e:
        print(e.read())


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("username")
    parser.add_argument("password")

    args = parser.parse_args()

    verification_token, session = get_csrf_data()
    if not verification_token:
        raise Exception("could not parse a request verification token from the page")

    post_login_form(args.username, args.password, verification_token, session)
