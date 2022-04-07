import pytz
from datetime import datetime
from ta.trend import SMAIndicator

from cryptellation.grapher import Grapher
from cryptellation.models import Period
from cryptellation.services import Candlesticks


class Graph(object):

    def __init__(self):
        self._grapher = Grapher()

    def run(self):
        start = datetime(2020, 7, 28).replace(tzinfo=pytz.utc)
        end = datetime(2020, 7, 29).replace(tzinfo=pytz.utc)

        data = Candlesticks().get('binance', 'BTC-USDC', Period.M5, start, end)
        self._grapher.candlesticks(data)

        sma = SMAIndicator(close=data['close'], window=10)
        self._grapher.line(sma.sma_indicator(), 'red')

        sma = SMAIndicator(close=data['close'], window=20)
        self._grapher.line(sma.sma_indicator(), 'yellow')

        sma = SMAIndicator(close=data['close'], window=40)
        self._grapher.line(sma.sma_indicator(), 'green')

        self._grapher.shade("2020-07-28T12:00:00Z", "2020-07-28T14:00:00Z",
                            "red", 0.1)
        self._grapher.show()


if __name__ == "__main__":
    Graph().run()