import collections
import uuid
from abc import ABCMeta, abstractmethod, ABC
from datetime import datetime
from typing import Deque

from logs import log
from savers.file_item import FileItem


class AbstractSaver(metaclass=ABCMeta):

    def __init__(self):
        self.file_queue: Deque[FileItem] = collections.deque()

        self.skipped = 0
        self.saved = 0

    def add(self, obj: str, key: str, datetime_obj: datetime, child: str, comment: str = None):
        fileitem = FileItem(datetime_obj, child, comment)
        log.info(f'Adding new file: {fileitem.base_name}')

        self.file_queue.append(fileitem)

    @abstractmethod
    def commit(self):
        """ process the file_queue and write the binary data """
        pass

    @abstractmethod
    def get_save_path(self, timestamp: datetime, file_name: str):
        """ get the target save path """
        pass

    @abstractmethod
    def exists(self, timestamp: datetime, file_name: str) -> bool:
        pass


class AbstractBucketSaver(AbstractSaver, ABC):
    """ base class for saving into cloud buckets """

    def __init__(self, bucket: str, access_id: str, secret_key: str):
        super().__init__()
        self.bucket = bucket
        self.access_id = access_id
        self.secret_key = secret_key

        self._test_connection()

    def _test_connection(self):
        pass
