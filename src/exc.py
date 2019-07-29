class NoEventsError(Exception):
    """ there was no response to the request query """


class NoMimeError(Exception):
    """ Can't save file with no mime type """


class NoTokenError(Exception):
    """ oauth token not provided"""
