from datetime import datetime, date
from json import JSONDecodeError

from pytz import timezone
from requests import RequestException

from exc import NoEventsError, UnauthorizedError
from logs import log
from savers import SAVER
from settings import client, conf
from utils import timestamp_to_date as t2d, to_daily_timestamp


def parse_events(start: [datetime, date], end: [datetime, date], num=300) -> (int, int, int):
    start = to_daily_timestamp(start)
    end = to_daily_timestamp(end)
    event_count = 0

    log.info(f'Request: {t2d(end)} - {t2d(start)}')

    params = {'direction': 'range',
              'earliest_event_time': start,
              'latest_event_time': end,
              'num_events': num}

    response = client.get(conf.EVENTS_URL, params=params)
    try:
        data = response.json()
        events = data['events']
        if not events and not conf.SKIP_NO_DATA_CHECK:
            raise NoEventsError('No events')
        for event in events:
            new_attachments = event.get('new_attachments', [])
            if not new_attachments:
                continue
            event_count += 1
            event_time = event['event_time']
            tz = event['tz']
            time = datetime.fromtimestamp(event_time, timezone(tz))
            obj_key = event['key']
            comment = event['comment']
            child_name = event.get('parent_member_display')

            # usually only one attachment, but just in case
            for att in new_attachments:
                SAVER.add(obj=obj_key, key=att['key'],
                          mime=att['mime_type'], timestamp=time, child=child_name, comment=comment)

        # finalize a batch of saver operations
        SAVER.commit()
        return event_count

    except RequestException:
        response.raise_for_status()
        return event_count

    except JSONDecodeError:
        if response.status_code == 401:
            raise UnauthorizedError('Auth token is invalid or expired')
        if not conf.SKIP_NO_DATA_CHECK:
            raise NoEventsError('No Events')
