import pytz
from datetime import datetime

from cryptellation.models import Period, Event
from cryptellation.backtester import Backtester


class Backtest(Backtester):

    def __init__(self, start_time: datetime, end_time: datetime):
        super().__init__(start_time=start_time, end_time=end_time)
        self.unique_order = False
        self.target_time = datetime(2020, 7, 28, 10,
                                    15).replace(tzinfo=pytz.utc)
        self.subscribe_ticks("binance", "BTC-USDC")

    def on_event(self, event: Event):
        if not self.unique_order and event.time() == self.target_time:
            self.order('market', 'binance', 'BTC-USDC', 'buy', 1)
            self.unique_order = True

    def on_end(self):
        self.order('market', 'binance', 'BTC-USDC', 'sell', 1)
        self.display('binance', 'BTC-USDC', Period.M1)


if __name__ == "__main__":
    b = Backtest(start_time=datetime(2020, 7, 28, 10).replace(tzinfo=pytz.utc),
                 end_time=datetime(2020, 7, 28, 12).replace(tzinfo=pytz.utc))
    b.run()
