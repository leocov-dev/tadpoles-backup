import calendar
from datetime import datetime, date

import filetype
from dateutil.relativedelta import relativedelta
from filetype.types import video

MAX_LINE_LEN = 80


def center_in_console(text):
    text_len = len(text) + 2
    if text_len >= MAX_LINE_LEN:
        return text[:MAX_LINE_LEN]

    pad = (MAX_LINE_LEN - text_len) // 2
    pad_str = '-' * pad

    return f'{pad_str} {text} {pad_str}'


def timestamp_to_date(a_timestamp: int) -> date:
    return datetime.fromtimestamp(a_timestamp).date()


def date_to_timestamp(a_date: date) -> int:
    return calendar.timegm(a_date.timetuple())


def to_timestamp_int(obj: [datetime, date]) -> int:
    if isinstance(obj, datetime):
        return int(obj.timestamp())
    elif isinstance(obj, date):
        return date_to_timestamp(obj)


BATCH_UNIT_MAP = {'days': 365,
                  'weeks': 52,
                  'months': 12}


def date_range_generator(batch_interval: int, batch_unit: str, start_date: [date, datetime], max_years: int):
    if isinstance(start_date, datetime):
        start_date = start_date.date()

    if batch_unit not in BATCH_UNIT_MAP:
        raise ValueError(f'batch_unit must be one of: {list(BATCH_UNIT_MAP.keys())}')

    current = start_date

    for _ in range(max_years * BATCH_UNIT_MAP[batch_unit]):
        previous = current - relativedelta(**{batch_unit: batch_interval})
        yield current, previous

        current = previous


class Mp4Compatible(video.Mp4):
    """
    More lenient mp4 detection for filetype package
    """

    def match(self, buf):
        if not self._is_isobmff(buf):
            return False

        major_brand, minor_version, compatible_brands = self._get_ftyp(buf)
        return any([cb in ['mp41', 'mp42'] for cb in compatible_brands])


filetype.add_type(Mp4Compatible())
