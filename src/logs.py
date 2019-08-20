import logging
import logging.config
from settings import conf


def config_logger(lvl='DEBUG'):
    logging.config.dictConfig({
        'version': 1,
        'formatters':
            {'default': {'format': '%(message)s'}},
        'handlers':
            {'console': {'level': lvl,
                         'class': 'logging.StreamHandler',
                         'formatter': 'default',
                         'stream': 'ext://sys.stdout'}},
        'loggers':
            {'tadpoles': {'level': lvl,
                          'handlers': ['console']}}
    })
    return logging.getLogger('tadpoles')


log = config_logger(conf.LOGGING_LEVEL)
