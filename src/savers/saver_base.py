import collections
from abc import ABCMeta, abstractmethod, ABC
from datetime import datetime
from typing import Deque

from logs import log
from savers.file_item import FileItem
from settings import concurrent_client, Config


class AbstractSaver(metaclass=ABCMeta):

    def __init__(self):
        self.file_queue: Deque[FileItem] = collections.deque()

        self.skipped = 0
        self.saved = 0

    def add(self, obj: str, key: str, datetime_obj: datetime, child: str, comment: str = None):
        try:
            future = concurrent_client.get(Config.ATTACHMENTS_URL, params={'obj': obj, 'key': key})
            file_item = FileItem(datetime_obj, child, comment, future)
            log.debug(f'Adding new file: {file_item.base_name}')

            self.file_queue.append(file_item)
        except TypeError as e:
            print(comment)
            raise e

    def commit(self):
        try:
            while self.file_queue:
                file_item = self.file_queue.popleft()
                if not file_item.ready:
                    self.file_queue.append(file_item)
                    continue
                if self._exists(file_item):
                    self.skipped += 1
                    continue
                self._write_file_item(file_item)
                self.saved += 1
        except Exception as e:
            log.exception(f'Commit Failure: {self.__class__.__name__}, {e}')

    @abstractmethod
    def _get_save_path(self, file_item: FileItem):
        """ get the target save path """

    @abstractmethod
    def _exists(self, file_item: FileItem) -> bool:
        """ does this file exist in the file system """

    @abstractmethod
    def _write_file_item(self, file_item):
        """ write the file item to the file system """


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
