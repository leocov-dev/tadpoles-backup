class TadpoleBackupError(Exception):
    """ base exception """


class NoEventsError(TadpoleBackupError):
    """ there was no response to the request query """


class NoMimeError(TadpoleBackupError):
    """ Can't save file with no mime type """


class NoTokenError(TadpoleBackupError):
    """ oauth token not provided"""


class SaveError(TadpoleBackupError):
    """ exception saving file """
