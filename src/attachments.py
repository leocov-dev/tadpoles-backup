import io

from requests import RequestException
from requests_toolbelt.downloadutils import stream

from settings import TadpolesConfig, client


def get_attachment(obj, key) -> bytes:
    response = client.get(TadpolesConfig.ATTACHMENTS_URL, params={'obj': obj, 'key': key})
    try:
        return response.content
    except RequestException:
        response.raise_for_status()


def stream_attachment(obj, key, bytes_obj: io.BytesIO):
    r = client.get(TadpolesConfig.ATTACHMENTS_URL, params={'obj': obj, 'key': key}, stream=True)
    stream.stream_response_to_file(r, path=bytes_obj)
