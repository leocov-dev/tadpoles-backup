import datetime

from dateutil.relativedelta import relativedelta

from events import parse_events
from exc import NoEventsError
from logs import log
from settings import conf


def main():
    previous = datetime.datetime.today()
    previous = previous.replace(hour=0, minute=0, second=0, microsecond=0)

    try:
        for _ in range(conf.MAX_YEARS * 12):
            earliest = previous - relativedelta(weeks=1)  # TODO: reset this to months=1
            log.info(f'Request: {previous.date()} - {earliest.date()}')
            parse_events(end=previous, start=earliest)
            previous = earliest
            break  # TODO: remove this
    except NoEventsError as e:
        log.info(f'Done with: {e.__class__.__name__}')
        return 0
    except Exception as e:
        log.exception(e)
        return 1

    log.info(f'Done with: MAX_YEARS ({conf.MAX_YEARS}) reached.')
    return 0


if __name__ == '__main__':
    exit(main())
