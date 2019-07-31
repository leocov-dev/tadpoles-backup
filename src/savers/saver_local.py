import os

from attachments import get_attachment
from savers.saver_base import AbstractSaver
from settings import conf


class LocalSaver(AbstractSaver):
    """ save to local directory """

    def __init__(self, target_dir=None):
        super().__init__()
        if not target_dir:
            target_dir = conf.LOCAL_TARGET_DIR
        self.target_dir = target_dir

    def save_path(self, child, file_name):
        return os.path.join(self.target_dir, child, file_name)

    def exists(self, child, file_name):
        return os.path.exists(self.save_path(child, file_name))

    def add(self, obj, key, ext, timestamp, child, comment=None):
        file_name = self.build_file_name(ext, timestamp, child, comment)
        if self.exists(child, file_name):
            return

        data = get_attachment(obj, key)
        self.file_queue.append((file_name, data))

    def commit(self):
        while self.file_queue:
            file_name, data = self.file_queue.pop()
