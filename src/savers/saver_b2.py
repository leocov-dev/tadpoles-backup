from savers.saver_base import AbstractBucketSaver


class B2Saver(AbstractBucketSaver):

    def __init__(self, bucket, access_id, secret_key):
        super().__init__(bucket, access_id, secret_key)

    def add(self, *args, **kwargs):
        pass

    def commit(self):
        pass

    def save_path(self, file_name):
        pass

    def exists(self, file_name):
        pass
