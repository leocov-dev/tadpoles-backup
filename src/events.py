from collections import defaultdict
from datetime import datetime, date
from json import JSONDecodeError

from pytz import timezone
from requests import RequestException

from exc import NoEventsError, UnauthorizedError
from logs import log
from savers import SAVER
from settings import client, conf, STR_DATE_FMT
from utils import timestamp_to_date as t2d, to_timestamp_int


def iter_events(start: [datetime, date], end: [datetime, date], num_events=5) -> (int, int, int):
    start = to_timestamp_int(start)
    end = to_timestamp_int(end)
    event_count = 0

    log.info(f'Request: {t2d(end)} - {t2d(start)}')

    params = {'direction': 'range',
              'earliest_event_time': start,
              'latest_event_time': end,
              'num_events': num_events}

    response = client.get(conf.EVENTS_URL, params=params)
    try:
        data = response.json()
        cursor = data['cursor']
        events = data['events']
        event_count += len(events)
        parse_events(events)
        print(cursor)
        while cursor:
            params['cursor'] = cursor
            response = client.get(conf.EVENTS_URL, params=params)
            data = response.json()
            cursor = data['cursor']
            events = data['events']
            event_count += len(events)
            parse_events(events)
            print(cursor)
        return event_count

    except RequestException:
        response.raise_for_status()
        return event_count

    except JSONDecodeError:
        if response.status_code == 401:
            raise UnauthorizedError('Auth token is invalid or expired')
        if not conf.SKIP_NO_DATA_CHECK:
            raise NoEventsError('No Events')


__no_comment_counter = defaultdict(lambda: defaultdict(int))


def parse_events(events):
    if not events and not conf.SKIP_NO_DATA_CHECK:
        raise NoEventsError('No events')

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
            comment = __no_comment_counter[child_name][time.date().strftime(STR_DATE_FMT)]

        # usually only one attachment, but just in case
        for att in new_attachments:
            SAVER.add(obj=obj_key, key=att['key'], datetime_obj=time, child=child_name, comment=comment)
