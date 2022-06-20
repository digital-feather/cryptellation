import grpc
from datetime import datetime
import threading
import queue
from typing import Dict, List
import iso8601
import json

from cryptellation.config import Config
from cryptellation.models import Account, Backtest

import cryptellation._genproto.backtests_pb2 as backtests
import cryptellation._genproto.backtests_pb2_grpc as backtests_grpc

class Backtests(object):

    def __init__(self):
        self._config = Config()
        self._channel = grpc.insecure_channel(self._config.backtests_url)
        self._stub = backtests_grpc.BacktestsServiceStub(self._channel)

    def create(self, start: datetime, end: datetime,
                        accounts: dict) -> Backtest:
        if start.tzinfo is None or start.tzinfo.utcoffset(start) is None:
            raise Exception("no timezone specified on start")

        id = self._stub.CreateBacktest(
            backtests.CreateBacktestRequest(
                start_time=start.isoformat(),
                end_time=end.isoformat(),
                accounts=self._account_to_grpc(accounts),
            )).id

        return Backtest(id)

    def _account_to_grpc(self, accounts: Dict[str, Account]):
        req_accounts = {}
        for exch, account in accounts.items():
            assets = {}
            for asset, quantity in account.assets.items():
                assets[asset] = quantity
            req_accounts[exch] = backtests.Account(assets=assets)
        return req_accounts
