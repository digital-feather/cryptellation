from datetime import datetime
from typing import List

from cryptellation.models.account import Account
from cryptellation.services.backtests import Backtests

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
                "binance": Account({"BTC": 1})
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
        self._id = self._backtests.create_backtest(
            start_time=self._config[Config.START_TIME]
        )
        self._events = self._backtests.listen_events(self._id)
        self.init()

    def init(self):
        pass

    def on_event(self, event):
        pass

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
                    break
                else:
                    self.on_event(e)
