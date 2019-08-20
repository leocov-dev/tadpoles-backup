import argparse
from datetime import datetime

from events import iter_events
from exc import NoEventsError, UnauthorizedError
from logs import log
from savers import SAVER
from savers.file_item import debug_file_type_info
from settings import conf, update_conf
from utils import date_range_generator, center_in_console

DELTA_MAP = {'days': 365,
             'weeks': 52,
             'months': 12}


def main(delta=1, delta_key='weeks', start_from=None):
    """ process events in blocks
    Examples:
        delta=1, delta_key='months'
        process events in 1 month chunks
    """
    if not start_from:
        start_from = datetime.now()
    total_events = 0

    return_code = 0
    log_func = log.info
    # default msg
    msg = f'Stopped after {conf.MAX_YEARS} years'

    try:
        log.info(center_in_console('Tadpoles-Backup'))
        log.info(f'Processing events in {delta} {delta_key.rstrip("s")} batches.')
        for current, previous in date_range_generator(delta, delta_key, start_from, conf.MAX_YEARS):
            event_count = iter_events(previous, current)
            total_events += event_count

    except NoEventsError:
        msg = 'Event block contained no events, exiting...'
    except UnauthorizedError as e:
        msg = f'{e.__class__.__name__}: {e}'
        log_func = log.error
        return_code = 1
    except Exception as e:
        log.exception(str(e))
        msg = f'Unexpected failure: {e.__class__.__name__}, {e}'
        log_func = log.error
        return_code = 1
    finally:
        log_func(f'{msg}')

        SAVER.commit()

        log.info(center_in_console(f'Checked: {total_events}, Skipped: {SAVER.skipped}, Saved {SAVER.saved}'))

        debug_file_type_info()
        return return_code


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Backup images and video from Tadpoles.com')
    parser.add_argument('--auth-token', help='Authentication Token')
    parser.add_argument('--start-date', help='Date in format YYYY-MM-DD')
    parser.add_argument('--max-years', default=6, type=int, help='Number of years past to query')
    parser.add_argument('--batch-unit', default='months', choices=['days', 'weeks', 'months'],
                        help='The unit for the batch interval.')
    parser.add_argument('--batch-interval', default=1, type=int,
                        help='The number of batch-units in each batch: 1 weeks, 2 months, etc.')
    subparser = parser.add_subparsers(title='modes', description='Saver Types')

    parser_local = subparser.add_parser('local')
    parser_local.add_argument('--save-path', help='Destination directory for files')

    parser_s3 = subparser.add_parser('s3')

    parser_b2 = subparser.add_parser('b2')

    args = parser.parse_args()

    update_conf(**vars(args))
    start_date = datetime.strptime(args.start_date, '%Y-%m-"d') if args.start_date else None

    exit(main(delta=args.batch_interval, delta_key=args.batch_unit, start_from=start_date))
