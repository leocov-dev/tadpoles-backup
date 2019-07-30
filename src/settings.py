import environs
import requests

from exc import NoTokenError

env = environs.Env()
env.read_env()


class Config:
    OAUTH_TOKEN = env('OAUTH_TOKEN')
    if not OAUTH_TOKEN:
        raise NoTokenError

    MAX_YEARS = env.int('MAX_YEARS', 10)
    MAX_FILE_NAME_LEN = 80
    API_URL = 'https://www.tadpoles.com/remote/v1'
    EVENTS_URL = f'{API_URL}/events'
    ATTACHMENTS_URL = f'{API_URL}/obj_attachment'
    SKIP_NO_DATA_CHECK = env.bool('SKIP_NO_DATA_CHECK', False)

    LOCAL_TARGET_DIR = env('LOCAL_TARGET_DIR')


def get_client():
    rc = requests.Session()
    rc.headers = {'Cookie': f'DgU00={conf.OAUTH_TOKEN}'}
    return rc


conf = Config()
client = get_client()
