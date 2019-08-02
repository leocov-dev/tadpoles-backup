import os
from datetime import datetime

from attachments import get_attachment
from savers.saver_base import AbstractSaver, FileItem
from settings import conf
from logs import log


class LocalSaver(AbstractSaver):
    """ save to a local directory """

    def __init__(self, target_dir: str = None):
        super().__init__()
        if not target_dir:
            target_dir = conf.LOCAL_TARGET_DIR
        self.target_dir = target_dir

    def save_path(self, timestamp: datetime, file_name: str):
        year_dir = os.path.join(self.target_dir, str(timestamp.year))
        os.makedirs(year_dir, exist_ok=True)
        return os.path.join(year_dir, file_name)

    def exists(self, timestamp: datetime, file_name: str):
        return os.path.exists(self.save_path(timestamp, file_name))

    def add(self, obj: str, key: str, mime: str, timestamp: datetime, child: str, comment: str = None):
        filename = self.build_file_name(mime, timestamp, child, comment)
        if self.exists(timestamp, filename):
            return

        log.info(f'Adding new file: {filename}')
        data = get_attachment(obj, key)
        self.file_queue.append(FileItem(data=data, mime=mime, timestamp=timestamp, filename=filename))

    def commit(self):
        while self.file_queue:
            file_item = self.file_queue.pop()
            with open(self.save_path(file_item.timestamp, file_item.filename), 'wb') as f:
                f.write(file_item.data)
