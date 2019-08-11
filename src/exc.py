class TadpoleBackupError(Exception):
    """ base exception """


class NoEventsError(TadpoleBackupError):
    """ there was no response to the request query """


class NoMimeError(TadpoleBackupError):
    """ Can't save file with no mime type """


class NoTokenError(TadpoleBackupError):
    """ oauth token not provided"""

    def __init__(self):
        super().__init__('Please provide an authentication token via the command line arguments or the .env file. '
                         'See README.md for more information.')


class SaveError(TadpoleBackupError):
    """ exception saving file """


class NoDataError(TadpoleBackupError):
    """ no data on FileItem """

    def __init__(self, fileitem):
        super().__init__(f'FileItem had no data: {fileitem}')


class UnauthorizedError(TadpoleBackupError):
    """ 401 error """
