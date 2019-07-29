import datetime

from dateutil.relativedelta import relativedelta

from events import parse_events
from exc import NoEventsError
from settings import conf


def main():
    previous = datetime.datetime.today()
    previous = previous.replace(hour=0, minute=0, second=0, microsecond=0)

    try:
        for _ in range(conf.MAX_YEARS * 12):
            earliest = previous - relativedelta(months=1)
            print(f'Request: {previous.date()} - {earliest.date()}')
            parse_events(end=previous, start=earliest)
            previous = earliest
    except NoEventsError as e:
        print(f'Done with: {e.__class__.__name__}')
        return 0

    print(f'Done with: MAX_YEARS ({conf.MAX_YEARS}) reached.')
    return 0


if __name__ == '__main__':
    exit(main())
