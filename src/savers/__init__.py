from exc import SaveError
from settings import conf


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


saver = get_saver()
