from typing import Dict

from .account import Account
from .candlesticks import Candlesticks

from cryptellation.config import Config

import cryptellation._genproto.livetests_pb2 as livetests
import cryptellation._genproto.livetests_pb2_grpc as livetests_grpc

class Livetest(object):
    def __init__(
        self,
        accounts: Dict[str, Account] = {
            "binance": Account({"USDC": 100000}),
        },
    ):
        self._config = Config()
        self._candlesticks = Candlesticks()
        self._channel = grpc.insecure_channel(self._config.livetests_url)
        self._stub = livetests_grpc.LivetestsServiceStub(self._channel)
        self._id = self._create(accounts)

    def _create(self, accounts) -> int:
        return self._stub.CreateLivetest(
            livetests.CreateLivetestRequest(
                accounts=self._account_to_grpc(accounts),
            )
        ).id

if __name__ == "__main__":
    t = Livetest()
