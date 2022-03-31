from typing import List


class Config():
    EXCHANGES_URL = "exchanges_url"
    CANDLESTICKS_URL = "candlesticks_url"
    BACKTESTS_URL = "backtests_url"

    def __init__(self, config: dict = {}):
        self._config = config
        self._check_default_config()

    def _check_default_config(self):
        if Config.EXCHANGES_URL not in self._config:
            self._config[Config.EXCHANGES_URL] = "127.0.0.1:9002"
        if Config.CANDLESTICKS_URL not in self._config:
            self._config[Config.CANDLESTICKS_URL] = "127.0.0.1:9003"
        if Config.BACKTESTS_URL not in self._config:
            self._config[Config.BACKTESTS_URL] = "127.0.0.1:9004"

    def keys() -> List[str]:
        return [
            Config.EXCHANGES_URL,
            Config.CANDLESTICKS_URL,
            Config.BACKTESTS_URL,
        ]

    def __getitem__(self, key):
        return self._config[key]
