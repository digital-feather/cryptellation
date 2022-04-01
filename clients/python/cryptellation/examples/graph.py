import pytz
from datetime import datetime

from cryptellation.grapher import Grapher
from cryptellation.models.period import Period


class Graph(object):

    def __init__(self):
        self._grapher = Grapher()

    def run(self):
        start = datetime(2020, 7, 28).replace(tzinfo=pytz.utc)
        end = datetime(2020, 7, 29).replace(tzinfo=pytz.utc)
        self._grapher.candlesticks('binance', 'BTC-USDC', Period.M5, start,
                                   end)
        self._grapher.simple_moving_average(10, 'red')
        self._grapher.simple_moving_average(20, 'yellow')
        self._grapher.simple_moving_average(50, 'green')
        # self._grapher.show()
        self._grapher.save()


if __name__ == "__main__":
    Graph().run()