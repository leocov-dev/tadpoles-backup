import environs
import requests

env = environs.Env()
env.read_env()


class Config:
    MAX_YEARS = env.int('MAX_YEARS', 10)
    MAX_FILE_NAME_LEN = 80
    OAUTH_TOKEN = env('OAUTH_TOKEN')
    API_URL = 'https://www.tadpoles.com/remote/v1'
    EVENTS_URL = f'{API_URL}/events'
    ATTACHMENTS_URL = f'{API_URL}/obj_attachment'
    SKIP_NO_DATA_CHECK = env.bool('SKIP_NO_DATA_CHECK', False)


def get_client():
    rc = requests.Session()
    rc.headers = {'Cookie': f'DgU00={conf.OAUTH_TOKEN}'}
    return rc


conf = Config()
client = get_client()
