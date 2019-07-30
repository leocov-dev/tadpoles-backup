from abc import ABCMeta, abstractmethod


class AbstractSaver(metaclass=ABCMeta):

    @abstractmethod
    def add(self, *args, **kwargs):
        pass

    @abstractmethod
    def commit(self):
        pass
