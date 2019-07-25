from settings import client, conf

URL = "https://www.tadpoles.com"

EVENTS_URL = f"{conf.API_URL}/events"


def parse_events():
    response = client.get(EVENTS_URL)
    data = response.json()
    events = data['events']
    for event in events:
        print(event.get('member_display'))
        entries = event.get('entries', [])
        for e in entries:
            note = e.get('note')
            attachment = e.get('attachment')
            if attachment:
                print(f'{note}: {attachment.get("filename")}')
