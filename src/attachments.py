from requests import RequestException

from settings import conf, client


def get_attachment(obj, key):
    response = client.get(conf.ATTACHMENTS_URL, params={'obj': obj, 'key': key})
    try:
        return response.content
    except RequestException:
        response.raise_for_status()
