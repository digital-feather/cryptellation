import pytz
from datetime import datetime

from cryptellation.backtester import Backtester, Config


class Visualizer(Backtester):

    def on_init(self):
        self.subscribe_ticks("binance", "BTC-USDC")
        self.subscribe_ticks("binance", "ETH-USDC")
        self.subscribe_ticks("binance", "LTC-USDC")

    def on_event(self, event):
        print(event.time, event.content)

    def on_exit(self):
        self.visual_summary('binance', 'BTC-USDC')


if __name__ == "__main__":
    config = Config({
        Config.START_TIME:
        datetime(2020, 7, 28).replace(tzinfo=pytz.utc),
        Config.END_TIME:
        datetime(2020, 7, 29).replace(tzinfo=pytz.utc),
    })
    Visualizer(config).run()