from savers.saver_base import AbstractBucketSaver


class S3Saver(AbstractBucketSaver):

    def __init__(self, bucket, access_id, secret_key):
        super().__init__(bucket, access_id, secret_key)

    def add(self, *args, **kwargs):
        pass

    def commit(self):
        pass

    def get_save_path(self, file_name):
        pass

    def exists(self, file_name):
        pass
