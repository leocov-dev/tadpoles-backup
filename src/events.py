import datetime
import re
import uuid

from pytz import timezone
from requests import RequestException

from attachments import get_attachment
from exc import NoEventsError
from settings import client, conf, saver


def parse_events(start, end, num=300):
    if isinstance(start, datetime.datetime):
        start = int(start.timestamp())
    if isinstance(end, datetime.datetime):
        end = int(end.timestamp())

    params = {'direction': 'range',
              'earliest_event_time': start,
              'latest_event_time': end,
              'num_events': num}

    response = client.get(conf.EVENTS_URL, params=params)
    try:
        data = response.json()
        events = data['events']
        if not events and not conf.SKIP_NO_DATA_CHECK:
            raise NoEventsError
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
                # todo: parse mimetype from attachment
                ext = "XXX"
                saver.add(ext=ext, timestamp=time, child=child_name, comment=comment)

        # finalize a batch of saver operations
        saver.commit()

        print(f'Got: {event_count} events.')
    except RequestException:
        response.raise_for_status()
