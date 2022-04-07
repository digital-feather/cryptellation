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
        self._id = self._backtests.create_backtest(start=start_time,
                                                   end=end_time,
                                                   accounts=accounts)
        self._actual_time = self._start_time
        self._events = self._backtests.listen_events(self._id)

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

    def subscribe_ticks(self, exchange_name, pair_symbol):
        self._backtests.subscribe_ticks(self._id, exchange_name, pair_symbol)

    def run(self):
        while True:
            if self._backtests.advance_backtest(self._id):
                break

            while True:
                event = self._events.get()

                if event.type == "end":
                    self._actual_time = event.time
                    break

                self.on_event(event)

        self.on_end()

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
        self._backtests.new_order(self._id, type, exchange, pair, side,
                                  quantity)

    def accounts(self):
        return self._backtests.accounts(self._id)

    def orders(self):
        return self._backtests.orders(self._id)