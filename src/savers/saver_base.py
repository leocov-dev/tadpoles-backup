import uuid
from abc import ABCMeta, abstractmethod, ABC

from settings import conf


class AbstractSaver(metaclass=ABCMeta):

    def __init__(self):
        self.file_queue = []

    @staticmethod
    def build_file_name(ext, timestamp, child_name, comment=None):
        if not comment:
            comment = str(uuid.uuid4()).split('-')[0]
        base_name = f'{timestamp.date().strftime("%Y.%m.%d")}-{child_name}-{comment}'
        max_name_len = conf.MAX_FILE_NAME_LEN - len(ext)
        file_name = f'{base_name[:max_name_len].rstrip("_")}{ext}'
        return file_name

    @abstractmethod
    def add(self, obj, key, ext, timestamp, child, comment=None):
        pass

    @abstractmethod
    def commit(self):
        pass

    @abstractmethod
    def save_path(self, child, file_name):
        pass

    @abstractmethod
    def exists(self, child, file_name):
        pass


class AbstractBucketSaver(AbstractSaver, ABC):
    """"""

    def __init__(self, bucket, access_id, secret_key):
        super().__init__()
        self.bucket = bucket
        self.access_id = access_id
        self.secret_key = secret_key

        self._test_connection()

    def _test_connection(self):
        pass
