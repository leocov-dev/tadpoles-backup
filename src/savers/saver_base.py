import collections
import io
import mimetypes
import re
import uuid
from abc import ABCMeta, abstractmethod, ABC
from datetime import datetime
from typing import Deque

import filetype

from attachments import get_attachment
from logs import log
from settings import conf

file_types = collections.defaultdict(int)


def print_file_type_info():
    log.debug(f'found types: {dict(file_types)}')


class FileItem:
    def __init__(self, data, mime, timestamp, filename, child, comment):
        self.data = data
        self.mime = mime
        self.timestamp = timestamp
        self.filename = filename
        self.child = child
        self.comment = comment.replace('\n', '')

        self._apply_metadata()

    def _apply_metadata(self):
        # out_data = io.BytesIO(self.data)
        # TODO: figuring this out
        img_type = filetype.guess(self.data)
        if not img_type:
            # for future development
            log.debug(self.data[:30])
            file_types[self.mime] += 1
        else:
            file_types[img_type.mime] += 1
        # if 'image' in self.mime:
        #     zeroth_ifd = {piexif.ImageIFD.Make: u"tadpoles-backup",
        #                   piexif.ImageIFD.Software: 'Python',
        #                   piexif.ImageIFD.ImageDescription: 'Josephine'
        #                   }
        #     exif_ifd = {piexif.ExifIFD.DateTimeOriginal: self.timestamp.strftime('%Y:%m:%d %H:%M:%S'),
        #                 piexif.ExifIFD.UserComment: UserComment.dump(self.comment)
        #                 }
        #
        #     exif_dict = {"0th": zeroth_ifd, "Exif": exif_ifd}
        #
        #     piexif.insert(piexif.dump(exif_dict), self.data, out_data)

        # self.data = out_data.read()


class AbstractSaver(metaclass=ABCMeta):

    def __init__(self):
        self.file_queue: Deque[FileItem] = collections.deque()

        self.skipped = 0
        self.saved = 0

    @staticmethod
    def _build_file_name(mime: str, timestamp: datetime, child_name: str, comment: str) -> str:
        comment = re.sub('\W+', '_', comment)
        comment = comment.rstrip('_')

        ext = mimetypes.guess_extension(mime)
        if ext in conf.REMAP_EXT:
            ext = conf.REMAP_EXT[ext]

        base_name = f'{timestamp.date().strftime("%Y.%m.%d")}-{child_name}-{comment}'
        max_name_len = conf.MAX_FILE_NAME_LEN - len(ext)
        file_name = f'{base_name[:max_name_len].rstrip("_")}{ext}'
        return file_name

    def add(self, obj: str, key: str, mime: str, timestamp: datetime, child: str, comment: str = None):
        *_, default_comment = str(uuid.uuid4).split('-')
        if not comment or comment == 'None':
            comment = default_comment

        filename = self._build_file_name(mime, timestamp, child, comment)
        if self.exists(timestamp, filename):
            self.skipped += 1
            return

        # TODO: temp
        # log.info(f'Adding new file: {filename}')

        self.file_queue.append(FileItem(data=get_attachment(obj, key),
                                        mime=mime,
                                        timestamp=timestamp,
                                        filename=filename,
                                        child=child,
                                        comment=comment))

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
