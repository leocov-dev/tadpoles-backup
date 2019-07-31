import environs
import requests

from exc import NoTokenError, SaveError

env = environs.Env()
env.read_env()


class Config:
    OAUTH_TOKEN = env('OAUTH_TOKEN')
    if not OAUTH_TOKEN:
        raise NoTokenError

    # sometimes mimetypes will not guess the 'commonly accepted' extension
    REMAP = {'.jpe': '.jpg'}

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


def get_saver():
    if conf.LOCAL_TARGET_DIR:
        from savers.saver_local import LocalSaver
        return LocalSaver(target_dir=conf.LOCAL_TARGET_DIR)

    if all([conf.S3_BUCKET_NAME, conf.S3_ACCESS_ID, conf.S3_SECRET_KEY]):
        from savers.saver_s3 import S3Saver
        return S3Saver(bucket=conf.S3_BUCKET_NAME, access_id=conf.S3_ACCESS_ID, secret_key=conf.S3_SECRET_KEY)

    if all([conf.B2_BUCKET_NAME, conf.B2_ACCOUNT_ID, conf.B2_ACCOUNT_KEY]):
        from savers.saver_b2 import B2Saver
        return B2Saver(bucket=conf.B2_BUCKET_NAME, access_id=conf.B2_ACCOUNT_ID, secret_key=conf.B2_ACCOUNT_KEY)

    raise SaveError('Could not get a saver with current environment variable settings.')


conf = Config()
client = get_client()
saver = get_saver()
