import collections
import mimetypes
import uuid
from abc import ABCMeta, abstractmethod, ABC
from datetime import datetime
from typing import Deque

from attachments import get_attachment
from logs import log
from settings import conf

FileItem = collections.namedtuple('FileItem', 'data mime timestamp filename')


class AbstractSaver(metaclass=ABCMeta):

    def __init__(self):
        self.file_queue: Deque[FileItem] = collections.deque()

        self.skipped = 0
        self.saved = 0

    @staticmethod
    def build_file_name(mime: str, timestamp: datetime, child_name: str, comment: str = None) -> str:
        ext = mimetypes.guess_extension(mime)
        if ext in conf.REMAP_EXT:
            ext = conf.REMAP_EXT[ext]
        if not comment:
            *_, comment = str(uuid.uuid4()).split('-')
        base_name = f'{timestamp.date().strftime("%Y.%m.%d")}-{child_name}-{comment}'
        max_name_len = conf.MAX_FILE_NAME_LEN - len(ext)
        file_name = f'{base_name[:max_name_len].rstrip("_")}{ext}'
        return file_name

    def add(self, obj: str, key: str, mime: str, timestamp: datetime, child: str, comment: str = None):
        filename = self.build_file_name(mime, timestamp, child, comment)
        if self.exists(timestamp, filename):
            self.skipped += 1
            return

        log.info(f'Adding new file: {filename}')
        data = get_attachment(obj, key)
        self.file_queue.append(FileItem(data=data, mime=mime, timestamp=timestamp, filename=filename))

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
