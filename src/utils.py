import calendar
from datetime import datetime, date

import filetype
from dateutil.relativedelta import relativedelta
from filetype.types.isobmff import IsoBmff


def timestamp_to_date(a_timestamp: int) -> date:
    return datetime.fromtimestamp(a_timestamp).date()


def date_to_timestamp(a_date: date) -> int:
    return calendar.timegm(a_date.timetuple())


def to_daily_timestamp(object: [datetime, date]) -> int:
    if isinstance(object, datetime):
        return int(object.timestamp())
    elif isinstance(object, date):
        return date_to_timestamp(object)


DELTA_MAP = {'days': 365,
             'weeks': 52,
             'months': 12}


def date_range_generator(delta: int, delta_key: str, start_date: [date, datetime], max_years: int):
    if isinstance(start_date, datetime):
        start_date = start_date.date()

    if delta_key not in DELTA_MAP:
        raise ValueError(f'delta_key must be one of: {list(DELTA_MAP.keys())}')

    current = start_date

    for _ in range(max_years * DELTA_MAP[delta_key]):
        previous = current - relativedelta(**{delta_key: delta})
        yield current, previous

        current = previous


class Mp4Lenient(IsoBmff):
    """
    More lenient mp4 detection for filetype package
    """
    MIME = 'video/mp4'
    EXTENSION = 'mp4'

    def __init__(self):
        super().__init__(
            mime=self.MIME,
            extension=self.EXTENSION
        )

    def match(self, buf):
        if not self._is_isobmff(buf):
            return False

        major_brand, minor_version, compatible_brands = self._get_ftyp(buf)
        return any([cb in ['mp41', 'mp42'] for cb in compatible_brands])


filetype.add_type(Mp4Lenient())
