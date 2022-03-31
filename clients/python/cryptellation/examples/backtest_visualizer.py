import pytz
from datetime import datetime

from cryptellation.models.period import Period
from cryptellation.backtester import Backtester, Config


class Visualizer(Backtester):

    def on_init(self):
        self.unique_order = False
        self.target_time = datetime(2020, 7, 28, 10,
                                    15).replace(tzinfo=pytz.utc)
        self.subscribe_ticks("binance", "BTC-USDC")

    def on_event(self, time: datetime, type: str, content: dict):
        if not self.unique_order and time == self.target_time:
            self.order('market', 'binance', 'BTC-USDC', 'buy', 1)
            self.unique_order = True

    def on_exit(self):
        self.order('market', 'binance', 'BTC-USDC', 'sell', 1)
        print(self.accounts())
        # self.visual_summary('binance', 'BTC-USDC', Period.M1)
        pass


if __name__ == "__main__":
    config = Config({
        Config.START_TIME:
        datetime(2020, 7, 28, 10).replace(tzinfo=pytz.utc),
        Config.END_TIME:
        datetime(2020, 7, 28, 12).replace(tzinfo=pytz.utc),
    })
    Visualizer(config).run()