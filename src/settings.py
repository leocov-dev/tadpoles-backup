import os

import environs
import requests
from requests_futures.sessions import FuturesSession

from exc import NoTokenError

env = environs.Env()
env.read_env()

STR_DATE_FMT = '%Y-%m-%d'


class Config:
    AUTH_TOKEN = env('OAUTH_TOKEN')
    if not AUTH_TOKEN:
        raise NoTokenError

    MAX_YEARS = env.int('MAX_YEARS', 10)
    EVENTS_PAGE_SIZE = env.int('EVENTS_PAGE_SIZE', 100)
    MAX_FILE_NAME_LEN = 80
    API_URL = 'https://www.tadpoles.com/remote/v1'
    EVENTS_URL = f'{API_URL}/events'
    ATTACHMENTS_URL = f'{API_URL}/obj_attachment'
    SKIP_NO_DATA_CHECK = env.bool('SKIP_NO_DATA_CHECK', False)
    LOGGING_LEVEL = env('LOGGING_LEVEL', 'INFO').upper()
    THREADED = env.bool('THREADED', False)

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

    @classmethod
    def update(cls, **kwargs):
        for k, v in kwargs.items():
            attr_name = k.upper()
            if hasattr(cls, attr_name):
                setattr(cls, attr_name, v)


def get_client(concurrent=False):
    if concurrent:
        cpu_count = os.cpu_count()
        session = FuturesSession(max_workers=cpu_count * 2 if cpu_count else 4)
    else:
        session = requests.Session()
    session.headers = {'Cookie': f'DgU00={Config.AUTH_TOKEN}'}
    return session


client = get_client(concurrent=False)
concurrent_client = get_client(concurrent=True)


