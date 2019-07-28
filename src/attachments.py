import mimetypes

from exc import NoMime
from settings import conf, client

REMAP = {'.jpe': '.jpg'}


def save_attachment(obj, key, base_file_name, saver, suffix=1):
    suffix = f"_{suffix}"
    response = client.get(conf.ATTACHMENTS_URL, params={'obj': obj, 'key': key})
    mime = response.headers.get('Content-Type')
    if not mime:
        raise NoMime

    ext = mimetypes.guess_extension(mime)
    if ext in REMAP:
        ext = REMAP[ext]

    max_name_len = conf.MAX_FILE_NAME_LEN - len(ext) - len(suffix)
    file_name = f'{base_file_name[:max_name_len].rstrip("_")}{suffix}{ext}'
    print(file_name)
