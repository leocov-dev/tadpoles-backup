from savers.saver_base import AbstractBucketSaver


class S3Saver(AbstractBucketSaver):
    def __init__(self, bucket, access_id, secret_key):
        super().__init__(bucket, access_id, secret_key)

    def _get_save_path(self, file_name):
        pass

    def _exists(self, file_name):
        pass

    def _write_file_item(self, file_item):
        pass

