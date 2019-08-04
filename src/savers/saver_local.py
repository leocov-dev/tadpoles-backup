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

    def get_save_path(self, timestamp: datetime, file_name: str) -> str:
        year_dir = os.path.join(self.target_dir, str(timestamp.year))
        os.makedirs(year_dir, exist_ok=True)
        return os.path.join(year_dir, file_name)

    def exists(self, timestamp: datetime, file_name: str) -> bool:
        return os.path.exists(self.get_save_path(timestamp, file_name))

    def commit(self):
        while self.file_queue:
            file_item = self.file_queue.popleft()
            with open(self.get_save_path(file_item.timestamp, file_item.filename), 'wb') as f:
                f.write(file_item.data)
            self.saved += 1
