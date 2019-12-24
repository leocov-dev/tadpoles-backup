import collections
import io
import re
import uuid
from concurrent.futures import Future
from datetime import datetime
from typing import Optional

import filetype
import piexif
from filetype.types import image, video
from piexif.helper import UserComment
from requests import Response

from exc import NoDataError
from logs import log
from settings import Config
from utils import Mp4Compatible

file_type_list = collections.defaultdict(int)


def debug_file_type_info():
    log.debug(f'found types: {dict(file_type_list)}')


class FileItem:
    def __init__(self, datetime_obj: datetime, child, comment, future: Future):
        future.add_done_callback(self._future_callback)
        self.future = future

        self.data: io.BytesIO = io.BytesIO()
        self.datetime = datetime_obj
        self.child = child
        if not comment:
            comment = str(uuid.uuid4().hex)
        self.comment = comment
        self.ext = 'bin'
        self.file_type: Optional[filetype.Type] = None
        self.exif = self._build_exif_dict()
        self.uid = uuid.uuid4()
        self._postprocessed = False

    @property
    def ready(self):
        return self.future.done() and self._postprocessed

    @property
    def year(self):
        return self.datetime.year

    @property
    def base_name(self):
        return self._build_file_name_base()

    @property
    def filename(self):
        if not self.file_type:
            self._guess_ext()
        return f'{self.base_name}.{self.ext}'

    @property
    def temp_filename(self):
        return f'{self.uid.hex}.{self.ext}'

    def __bytes__(self):
        return self.data.getvalue()

    def __str__(self):
        return f'{self.datetime}'

    def _future_callback(self, future: Future):
        if not future.cancelled():
            response: Response = future.result(timeout=3)
            self.data = io.BytesIO(response.content)
            self._guess_ext()
            self._apply_metadata()
            self._postprocessed = True

    def _guess_ext(self):
        if not self.data:
            raise NoDataError(self)

        self.file_type: filetype.Type = filetype.guess(self.data)
        if not self.file_type:
            log.debug(f'Could not determine file type from header[:262]: {self.data[:262]}')
            file_type_list['bin'] += 1
        else:
            self.ext = self.file_type.extension
            file_type_list[self.file_type.extension] += 1

    def _build_file_name_base(self) -> str:
        comment = re.sub('\W+', '_', self.comment)
        comment = comment.rstrip('_')

        base_name = f'{self.datetime.strftime("%Y.%m.%d")}-{self.child}-{comment}'.rstrip('_')
        return base_name[:Config.MAX_FILE_NAME_LEN]

    def _apply_metadata(self):
        if not self.data:
            raise NoDataError(self)

        try:
            if isinstance(self.file_type, image.Jpeg):
                self._jpg_metadata()
            elif isinstance(self.file_type, image.Png):
                print('>PNG')
                from PIL import Image
                im = Image.open(self.data)
                print(f'{im}')
                rgb = im.convert('RGB')
                print(f'{rgb}')
                rgb.save(self.data, format='jpeg')
                self._jpg_metadata()
                self.file_type = image.Jpeg()
                self._guess_ext()
                print(f'{self.ext}>')
            elif any([isinstance(self.file_type, t) for t in [Mp4Compatible, video.Mp4]]):
                self._mp4_metadata()
        except Exception as e:
            log.exception(e)
            exit(str(e))

    def _build_exif_dict(self):
        log.debug(f'Applying JPEG metadata')
        zeroth_ifd = {piexif.ImageIFD.Make: "tadpoles-backup",
                      piexif.ImageIFD.Software: 'Python',
                      piexif.ImageIFD.ImageDescription: self.child
                      }
        exif_ifd = {piexif.ExifIFD.DateTimeOriginal: self.datetime.strftime('%Y:%m:%d %H:%M:%S'),
                    piexif.ExifIFD.UserComment: UserComment.dump(self.comment)
                    }

        return {"0th": zeroth_ifd, "Exif": exif_ifd}

    def _jpg_metadata(self):
        log.debug('Applying JPEG metadata...')
        piexif.insert(piexif.dump(self.exif), self.data)

    def _mp4_metadata(self):
        log.debug(f'{self.ext} Not yet implemented...')
