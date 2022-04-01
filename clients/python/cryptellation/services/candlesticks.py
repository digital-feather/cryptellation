import grpc
from datetime import datetime, timedelta
import pytz
import pandas as pd

from cryptellation.config import Config
from cryptellation.models.period import Period

import cryptellation.services.genproto.candlesticks_pb2 as candlesticks
import cryptellation.services.genproto.candlesticks_pb2_grpc as candlesticks_grpc


class Candlesticks(object):

    def __init__(self):
        self._config = Config()
        self._channel = grpc.insecure_channel(
            self._config[Config.CANDLESTICKS_URL])
        self._stub = candlesticks_grpc.CandlesticksServiceStub(self._channel)

    def get(
        self,
        exchange: str,
        pair: str,
        period: Period,
        start: datetime,
        end: datetime,
        limit: int = 0,
    ):
        if start.tzinfo is None or start.tzinfo.utcoffset(start) is None:
            raise Exception("no timezone specified on start")
        if end.tzinfo is None or end.tzinfo.utcoffset(end) is None:
            raise Exception("no timezone specified on end")

        cs = self._stub.ReadCandlesticks(
            candlesticks.ReadCandlesticksRequest(
                exchange_name=exchange,
                pair_symbol=pair,
                period_symbol=str(period),
                start=start.isoformat(),
                end=end.isoformat(),
                limit=limit,
            )).candlesticks

        list = {}
        for c in cs:
            list[c.time] = [c.open, c.high, c.low, c.close, c.volume]

        df = pd.DataFrame.from_dict(
            list,
            orient='index',
            columns=['open', 'high', 'low', 'close', 'volume'])
        df.index.names = ['time']
        return df


if __name__ == "__main__":
    now = datetime.utcnow().replace(tzinfo=pytz.utc)
    start = now - timedelta(hours=0, minutes=5)
    end = now - timedelta(hours=0, minutes=3)
    cs = Candlesticks().get("binance", "BTC-USDT", Period.M1, start, end, 0)
    print(cs)