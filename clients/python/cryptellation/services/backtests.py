import grpc
from datetime import datetime
import threading
import queue

from cryptellation.config import Config
from cryptellation.models.account import Account

import cryptellation.services.genproto.backtests_pb2 as backtests
import cryptellation.services.genproto.backtests_pb2_grpc as backtests_grpc


class BacktestsEventsQueue(threading.Thread):

    def __init__(self, generator):
        threading.Thread.__init__(self)
        self._generator = generator
        self._events_queue = queue.Queue(maxsize=0)

    def run(self):
        for event in self._generator:
            self._events_queue.put(event)

    def get(self):
        return self._events_queue.get()


class Backtests(object):

    def __init__(self):
        self._config = Config()
        self._channel = grpc.insecure_channel(
            self._config[Config.BACKTESTS_URL])
        self._stub = backtests_grpc.BacktestsServiceStub(self._channel)

    def create_backtest(self, start: datetime, end: datetime, accounts: dict):
        if start.tzinfo is None or start.tzinfo.utcoffset(start) is None:
            raise Exception("no timezone specified on start")

        return self._stub.CreateBacktest(
            backtests.CreateBacktestRequest(
                start_time=start.isoformat(),
                end_time=end.isoformat(),
                accounts=self._account_to_grpc(accounts),
            )).id

    def _account_to_grpc(self, accounts: dict):
        req_accounts = {}
        for exch, account in accounts.items():
            assets = {}
            for asset, quantity in account.assets.items():
                assets[asset] = quantity
            req_accounts[exch] = backtests.Account(assets=assets)
        return req_accounts

    def advance_backtest(self, id):
        return self._stub.AdvanceBacktest(
            backtests.AdvanceBacktestRequest(id=id, )).finished

    def subscribe_ticks(self, id, exchange_name, pair_symbol):
        self._stub.SubscribeToBacktestEvents(
            backtests.SubscribeToBacktestEventsRequest(
                id=id,
                exchange_name=exchange_name,
                pair_symbol=pair_symbol,
            ))

    def listen_events(self, id):
        req = backtests.ListenBacktestRequest(id=id)
        q = BacktestsEventsQueue(self._stub.ListenBacktest(req))
        q.start()
        return q

    def new_order(self, id: int, type: str, exchange: str, pair: str,
                  side: str, quantity: float):
        req = backtests.CreateBacktestOrderRequest(
            backtest_id=id,
            type=type,
            exchange_name=exchange,
            pair_symbol=pair,
            side=side,
            quantity=quantity,
        )
        self._stub.CreateBacktestOrder(req)

    def accounts(self, id: int):
        req = backtests.AccountsRequest(backtest_id=id, )
        resp = self._stub.Accounts(req)
        return self._grpc_to_accounts(resp.accounts)

    def _grpc_to_accounts(self, resp: dict):
        accounts = {}
        for exch, account in resp.items():
            assets = {}
            for asset, quantity in account.assets.items():
                assets[asset] = quantity
            accounts[exch] = Account(assets)
        return accounts