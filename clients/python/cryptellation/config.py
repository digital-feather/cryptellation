from typing import List


class Config:
    def __init__(
        self,
        exchanges_url: str = "127.0.0.1:9002",
        candlesticks_url: str = "127.0.0.1:9003",
        backtests_url: str = "127.0.0.1:9004",
        ticks_url: str = "127.0.0.1:9005",
        livetests_url: str = "127.0.0.1:9006",
    ):
        self.exchanges_url = exchanges_url
        self.candlesticks_url = candlesticks_url
        self.backtests_url = backtests_url
        self.ticks_url = ticks_url
        self.livetests_url = livetests_url
