from datetime import datetime

from events import parse_events
from exc import NoEventsError, UnauthorizedError
from logs import log
from savers.saver_base import print_file_type_info
from settings import conf
from utils import date_range_generator

DELTA_MAP = {'days': 365,
             'weeks': 52,
             'months': 12}


def main(delta=1, delta_key='weeks', start_from=datetime.now()):
    """ process events in blocks
    Examples:
        delta=1, delta_key='months'
        process events in 1 month chunks
    """
    total_events = 0
    total_skipped = 0
    total_saved = 0

    # default msg
    msg = f'Stopped after {conf.MAX_YEARS} years'

    try:
        log.info(f'Processing events in {delta} {delta_key.rstrip("s")} batches.\n')
        for current, previous in date_range_generator(delta, delta_key, start_from, conf.MAX_YEARS):
            event_count, skipped, saved = parse_events(previous, current)
            total_events += event_count
            total_skipped += skipped
            total_saved += saved
            # break  # TODO: remove this
    except NoEventsError as e:
        msg = 'Event block contained no events, exiting...'
        log.warning(str(e))
        return
    except UnauthorizedError as e:
        log.critical(str(e))
        msg = e.__class__.__name__
        return 1
    except Exception as e:
        msg = 'Unexpected failure.'
        log.exception(e)
        msg = e.__class__.__name__
        return 1
    finally:
        log.info(f'\nDone with: {msg}')
        log.info(f'Checked: {total_events}, Skipped: {total_skipped}, Saved {total_saved}')

        print_file_type_info()


if __name__ == '__main__':
    exit(main(delta_key='months'))
