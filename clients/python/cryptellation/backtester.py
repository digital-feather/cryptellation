from datetime import datetime
from re import S
import plotly.graph_objects as go
from typing import List
import iso8601
import json

from cryptellation.models.period import Period
from cryptellation.models.account import Account
from cryptellation.services.backtests import Backtests
from cryptellation.services.candlesticks import Candlesticks
from cryptellation.grapher import Grapher


class Config(object):
    START_TIME = "start_time"
    END_TIME = "end_time"
    START_ACCOUNTS = "start_accounts"

    def __init__(self, config: dict = {}):
        self._config = config
        self._check_config()

    def _check_config(self):
        if Config.START_TIME not in self._config:
            raise ValueError('No start time specified')
        if type(self._config[Config.START_TIME]) is not datetime:
            raise ValueError('Invalid start time')

        if Config.START_ACCOUNTS not in self._config:
            self._config[Config.START_ACCOUNTS] = {
                "binance": Account({"USDC": 100000})
            }

        if Config.END_TIME not in self._config:
            self._config[Config.END_TIME] = datetime.now()

    def keys() -> List[str]:
        return [
            Config.START_TIME,
            Config.END_TIME,
            Config.START_ACCOUNTS,
        ]

    def __getitem__(self, key):
        return self._config[key]


class Backtester(object):

    def __init__(self, config: Config):
        self._config = config
        self._backtests = Backtests()
        self._candlesticks = Candlesticks()
        self._id = self._backtests.create_backtest(
            start=self._config[Config.START_TIME],
            end=self._config[Config.END_TIME],
            accounts=self._config[Config.START_ACCOUNTS])
        self._actual_time = self._config[Config.START_TIME]
        self._events = self._backtests.listen_events(self._id)
        self.on_init()

    def on_init(self):
        pass

    def on_event(self, time: datetime, type: str, content: dict):
        pass

    def on_exit(self):
        pass

    def display(self, exchange: str, pair: str, period: Period):
        p = Grapher()

        start = self._config[Config.START_TIME]
        end = self._config[Config.END_TIME]
        p.candlesticks(exchange, pair, period, start, end)

        p.orders(self.orders())

        p.show()

    def subscribe_ticks(self, exchange_name, pair_symbol):
        self._backtests.subscribe_ticks(self._id, exchange_name, pair_symbol)

    def run(self):
        while True:
            finished = self._backtests.advance_backtest(self._id)
            if finished:
                break

            while True:
                e = self._events.get()

                if e.type == "end":
                    self._actual_time = iso8601.parse_date(e.time)
                    break

                self.on_event(iso8601.parse_date(e.time), e.type,
                              json.loads(e.content))

        self.on_exit()

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
        return self._backtests.accounts(self._id).accounts

    def orders(self):
        return self._backtests.orders(self._id).orders