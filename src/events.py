import requests

URL = "https://www.tadpoles.com"

EVENTS_URL = "/remote/v1/events"


def parse_events():
    response = requests.get(URL + EVENTS_URL)

    print(response.headers)
    print(response.content)
