import os

from exc import SaveError
from savers.saver_base import AbstractSaver
from settings import conf


class LocalSaver(AbstractSaver):
    """ save to local directory """

    def __init__(self, target_dir=None):
        if not target_dir:
            target_dir = conf.LOCAL_TARGET_DIR
        self.target_dir = target_dir

    def add(self, file_name):
        expected_path = os.path.join(self.target_dir, file_name)
        if os.path.exists(expected_path):
            raise SaveError(f'Path already exists: {expected_path}')

    def commit(self):
        pass
