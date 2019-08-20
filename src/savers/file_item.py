import collections
import io
import re
import uuid
from datetime import datetime

import filetype
import piexif
from filetype.types import image, video
from piexif.helper import UserComment

from exc import NoDataError
from logs import log
from settings import conf
from utils import Mp4Compatible

file_type_list = collections.defaultdict(int)


def debug_file_type_info():
    log.debug(f'found types: {dict(file_type_list)}')


class FileItem:
    def __init__(self, datetime_obj: datetime, child, comment, data=None):
        if not data:
            data = io.BytesIO()
        self.data = data
        self.datetime = datetime_obj
        self.child = child
        if not comment:
            comment = str(uuid.uuid4().hex)
        self.comment = comment
        self.ext = 'bin'
        self.file_type = None
        self.uid = uuid.uuid4()

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

    def __str__(self):
        return f'{self.datetime}'

    def download_data(self):
        pass

    def _guess_ext(self):
        if not self.data:
            raise NoDataError(self)

        self.file_type = filetype.guess(self.data)
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
        return base_name[:conf.MAX_FILE_NAME_LEN]

    def _apply_metadata(self):
        if not self.data:
            raise NoDataError(self)

        if isinstance(self.file_type, image.Jpeg):
            self._jpg_metadata()
        elif isinstance(self.file_type, image.Png):
            from PIL import Image
            im = Image.open(self.data)
            rgb = im.convert('RGB')
            rgb.save(self.data, format='jpeg')
            self._jpg_metadata()
        elif any([isinstance(self.file_type, t) for t in [Mp4Compatible, video.Mp4]]):
            self._mp4_metadata()

    def _jpg_metadata(self):
        zeroth_ifd = {piexif.ImageIFD.Make: "tadpoles-backup",
                      piexif.ImageIFD.Software: 'Python',
                      piexif.ImageIFD.ImageDescription: self.child
                      }
        exif_ifd = {piexif.ExifIFD.DateTimeOriginal: self.datetime.strftime('%Y:%m:%d %H:%M:%S'),
                    piexif.ExifIFD.UserComment: UserComment.dump(self.comment)
                    }

        exif_dict = {"0th": zeroth_ifd, "Exif": exif_ifd}

        piexif.insert(piexif.dump(exif_dict), self.data)

    def _mp4_metadata(self):
        log.debug(f'{self.ext} Not yet implemented...')
