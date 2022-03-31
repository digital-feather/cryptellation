from datetime import datetime

from cryptellation.backtester import Backtester, Config


class Visualizer(Backtester):

    def init(self):
        self.subscribe_ticks("binance", "BTC-USDC")
        self.subscribe_ticks("binance", "ETH-USDC")
        self.subscribe_ticks("binance", "LTC-USDC")

    def on_event(self, event):
        print(event.time, event.content)


if __name__ == "__main__":
    config = Config({
        Config.START_TIME: datetime(2020, 6, 28),
        Config.END_TIME: datetime(2020, 6, 28, 2),
    })
    Visualizer(config).run()
    print("end")