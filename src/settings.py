import environs
import requests

from exc import NoTokenError

env = environs.Env()
env.read_env()


class Config:
    OAUTH_TOKEN = env('OAUTH_TOKEN')
    if not OAUTH_TOKEN:
        raise NoTokenError

    # sometimes mimetypes will not guess the 'commonly accepted' extension
    REMAP_EXT = {'.jpe': '.jpg'}

    MAX_YEARS = env.int('MAX_YEARS', 10)
    MAX_FILE_NAME_LEN = 80
    API_URL = 'https://www.tadpoles.com/remote/v1'
    EVENTS_URL = f'{API_URL}/events'
    ATTACHMENTS_URL = f'{API_URL}/obj_attachment'
    SKIP_NO_DATA_CHECK = env.bool('SKIP_NO_DATA_CHECK', False)

    # save files to a local directory
    LOCAL_TARGET_DIR = env('LOCAL_TARGET_DIR', None)

    # save files to amazon s3
    S3_BUCKET_NAME = env('S3_BUCKET_NAME', None)
    S3_ACCESS_ID = env('S3_ACCESS_ID', None)
    S3_SECRET_KEY = env('S3_SECRET_KEY', None)

    # save files to backblaze b2
    B2_BUCKET_NAME = env('B2_BUCKET_NAME', None)
    B2_ACCOUNT_ID = env('B2_ACCOUNT_ID', None)
    B2_ACCOUNT_KEY = env('B2_ACCOUNT_KEY', None)


def get_client():
    rc = requests.Session()
    rc.headers = {'Cookie': f'DgU00={conf.OAUTH_TOKEN}'}
    return rc


conf = Config()
client = get_client()
