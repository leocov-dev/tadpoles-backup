import os

from logs import log
from savers.file_item import FileItem
from savers.saver_base import AbstractSaver
from settings import Config


class LocalSaver(AbstractSaver):
    """ save to a local directory """

    def __init__(self, target_dir: str = None):
        super().__init__()
        if not target_dir:
            target_dir = Config.LOCAL_TARGET_DIR
        self.target_dir = target_dir

    def _get_save_path(self, file_item: FileItem) -> str:
        year_dir = os.path.join(self.target_dir, str(file_item.datetime.year))
        os.makedirs(year_dir, exist_ok=True)
        return os.path.join(year_dir, file_item.filename)

    def _exists(self, file_item: FileItem) -> bool:
        return os.path.exists(self._get_save_path(file_item))

    def _write_file_item(self, file_item):
        with open(self._get_save_path(file_item), 'wb') as f:
            f.write(bytes(file_item))
