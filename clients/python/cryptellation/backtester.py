from datetime import datetime
from typing import Dict, List

from cryptellation.models import Period, Account, Event
from cryptellation.services import Backtests, Candlesticks
from cryptellation.grapher import Grapher


class Backtester(object):

    def __init__(self,
                 start_time: datetime,
                 end_time: datetime = datetime.now(),
                 accounts: Dict[str, Account] = {
                     "binance": Account({"USDC": 100000}),
                 }):
        self._start_time = start_time
        self._end_time = end_time
        self._backtests = Backtests()
        self._candlesticks = Candlesticks()
        self._backtest = self._backtests.create(start=start_time,
                                                   end=end_time,
                                                   accounts=accounts)
        self._actual_time = self._start_time
        self._events = self._backtest.listen()

    def on_event(self, event: Event):
        pass

    def on_end(self):
        pass

    def display(self, exchange: str, pair: str, period: Period):
        p = Grapher()

        start = self._start_time
        end = self._end_time
        cs = self._candlesticks.get(exchange, pair, period, start, end)
        p.candlesticks(cs)

        p.orders(self.orders())

        p.show()

    def subscribe(self, exchange_name, pair_symbol):
        self._backtest.subscribe(exchange_name, pair_symbol)

    def actual_time(self) -> datetime:
        return self._actual_time

    def run(self):
        finished = False
        while finished is False:
            self._events.next()

            while True:
                event = self._events.get()

                if event.type == "status":
                    finished = event.content['finished']
                    self._actual_time = event.time
                    break

                self.on_event(event)

        return self.on_end()

    def candlesticks(
        self,
        exchange: str,
        pair: str,
        period: Period,
        relative_start: int,
        relative_end: int = 0,
        limit: int = 0,
    ):
        start = self._actual_time - relative_start * period.duration()
        end = self._actual_time - relative_end * period.duration()
        return self._candlesticks.get(exchange, pair, period, start, end,
                                      limit)

    def order(self, type: str, exchange: str, pair: str, side: str,
              quantity: float):
        self._backtest.order(type, exchange, pair, side,
                                  quantity)

    def accounts(self):
        return self._backtest.accounts()

    def orders(self):
        return self._backtest.orders()
