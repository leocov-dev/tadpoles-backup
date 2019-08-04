from datetime import datetime, date
import re
import uuid

from pytz import timezone
from requests import RequestException

from exc import NoEventsError
from logs import log
from savers import saver
from settings import client, conf
from utils import timestamp_to_date as t2d, date_to_timestamp, to_daily_timestamp


def parse_events(start: [datetime, date], end: [datetime, date], num=300) -> (int, int, int):
    start = to_daily_timestamp(start)
    end = to_daily_timestamp(end)

    log.info(f'Request: {t2d(end)} - {t2d(start)}')
    return 0, 0, 0

    params = {'direction': 'range',
              'earliest_event_time': start,
              'latest_event_time': end,
              'num_events': num}

    response = client.get(conf.EVENTS_URL, params=params)
    try:
        data = response.json()
        events = None  # data['events']
        if not events and not conf.SKIP_NO_DATA_CHECK:
            raise NoEventsError(
                f'Event-range: {t2d(end)} - {t2d(start)} returned no events')
        event_count = 0
        for event in events:
            new_attachments = event.get('new_attachments', [])
            if not new_attachments:
                continue
            event_count += 1
            event_time = event['event_time']
            tz = event['tz']
            time = datetime.datetime.fromtimestamp(event_time, timezone(tz))
            obj_key = event['key']
            default_comment = str(uuid.uuid4).split('-')[0]
            comment = event.get('comment', default_comment)
            if comment:
                if comment == 'None':
                    comment = default_comment
                else:
                    comment = re.sub('\W+', '_', comment)
                    comment = comment.rstrip('_')
            child_name = event.get('parent_member_display')

            # usually only one attachment, but just in case
            for att in new_attachments:
                saver.add(obj=obj_key, key=att['key'],
                          mime=att['mime_type'], timestamp=time, child=child_name, comment=comment)

        # finalize a batch of saver operations
        saver.commit()
        return event_count, saver.skipped, saver.saved

    except RequestException:
        response.raise_for_status()
