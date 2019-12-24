from collections import defaultdict
from datetime import datetime, date
from json import JSONDecodeError

from pytz import timezone
from requests import RequestException

from exc import NoEventsError, UnauthorizedError
from logs import log
from savers import SAVER
from settings import client, Config, STR_DATE_FMT
from utils import timestamp_to_date as t2d, to_timestamp_int


def iter_events(start: [datetime, date], end: [datetime, date], event_page_size=Config.EVENTS_PAGE_SIZE) -> (int, int, int):
    start = to_timestamp_int(start)
    end = to_timestamp_int(end)
    events_total = 0

    log.info(f'Request: {t2d(end)} - {t2d(start)}')

    params = {'direction': 'range',
              'earliest_event_time': start,
              'latest_event_time': end,
              'num_events': event_page_size}

    cursor = True
    while cursor:
        if isinstance(cursor, str):
            params['cursor'] = cursor
        response = client.get(Config.EVENTS_URL, params=params)
        response.raise_for_status()
        try:
            data = response.json()
            cursor = data['cursor']
            events = data['events']
            events_count = len(events)
            if events_count > 0:
                events_total += events_count
                parse_events(events)

        except JSONDecodeError as e:
            if response.status_code == 401:
                raise UnauthorizedError('Auth token is invalid or expired')
            if not Config.SKIP_NO_DATA_CHECK:
                raise NoEventsError('Could not parse JSON events')

    if events_total == 0 and not Config.SKIP_NO_DATA_CHECK:
        raise NoEventsError('No events in block.')

    return events_total


__no_comment_counter = defaultdict(lambda: defaultdict(int))


def parse_events(events):
    """ inspect the events in this batch and save any attachments that are found """
    for event in events:
        # looking for events with file attachments
        new_attachments = event.get('new_attachments', [])
        if not new_attachments:
            continue

        event_time = event['event_time']
        tz = event['tz']
        time = datetime.fromtimestamp(event_time, timezone(tz))
        obj_key = event['key']
        comment = None if event['comment'] == 'None' else event['comment']
        child_name = event.get('parent_member_display')
        if not comment:
            __no_comment_counter[child_name][time.date().strftime(STR_DATE_FMT)] += 1
            count = __no_comment_counter[child_name][time.date().strftime(STR_DATE_FMT)]
            comment = f"{count:04}"

        # usually only one attachment, but just in case
        for att in new_attachments:
            SAVER.add(obj=obj_key, key=att['key'], datetime_obj=time, child=child_name, comment=comment)
