import os

from logs import log
from savers.file_item import FileItem
from savers.saver_base import AbstractSaver
from settings import conf


class LocalSaver(AbstractSaver):
    """ save to a local directory """

    def __init__(self, target_dir: str = None):
        super().__init__()
        if not target_dir:
            target_dir = conf.LOCAL_TARGET_DIR
        self.target_dir = target_dir

    def get_save_path(self, file_item: FileItem) -> str:
        year_dir = os.path.join(self.target_dir, str(file_item.datetime.year))
        os.makedirs(year_dir, exist_ok=True)
        return os.path.join(year_dir, file_item.filename)

    def exists(self, file_item: FileItem) -> bool:
        return os.path.exists(self.get_save_path(file_item))

    def commit(self):
        log.info(f'Saving files in: {self.target_dir}...')
        try:
            while self.file_queue:
                file_item = self.file_queue.popleft()
                if self.exists(file_item):
                    self.skipped += 1
                    continue
                # with open(self.get_save_path(file_item), 'wb') as f:
                #     f.write(file_item.data)
                self.saved += 1
        except Exception as e:
            log.exception(f'Commit Failure: {self.__class__.__name__}, {e}')
