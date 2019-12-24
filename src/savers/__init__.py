from exc import SaveError
from settings import Config


def get_saver():
    if Config.LOCAL_TARGET_DIR:
        from savers.saver_local import LocalSaver
        return LocalSaver(target_dir=Config.LOCAL_TARGET_DIR)

    if all([Config.S3_BUCKET_NAME, Config.S3_ACCESS_ID, Config.S3_SECRET_KEY]):
        from savers.saver_s3 import S3Saver
        return S3Saver(bucket=Config.S3_BUCKET_NAME, access_id=Config.S3_ACCESS_ID, secret_key=Config.S3_SECRET_KEY)

    if all([Config.B2_BUCKET_NAME, Config.B2_ACCOUNT_ID, Config.B2_ACCOUNT_KEY]):
        from savers.saver_b2 import B2Saver
        return B2Saver(bucket=Config.B2_BUCKET_NAME, access_id=Config.B2_ACCOUNT_ID, secret_key=Config.B2_ACCOUNT_KEY)

    raise SaveError('Could not get a saver with current environment variable settings.')


SAVER = get_saver()
