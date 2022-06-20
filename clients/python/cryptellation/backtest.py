import grpc
import threading
import queue
from typing import Dict, List
import iso8601
import json
from datetime import datetime

from .account import Account
from .event import Event
from .order import Order
from .period import Period
from .grapher import Grapher
from .candlesticks import Candlesticks

from cryptellation.config import Config

import cryptellation._genproto.backtests_pb2 as backtests
import cryptellation._genproto.backtests_pb2_grpc as backtests_grpc


class BacktestEvents(threading.Thread):
    def __init__(self, id, stub: backtests_grpc.BacktestsServiceStub):
        threading.Thread.__init__(self)
        self._id = id
        self._next_queue = queue.Queue(maxsize=0)
        self._events_queue = queue.Queue(maxsize=0)
        queue_iterator = iter(self._next_queue.get, None)
        self._generator = stub.ListenBacktest(queue_iterator)
        self.start()

    def run(self):
        for event in self._generator:
            self._events_queue.put(event)

    def get(self) -> Event:
        e = self._events_queue.get()
        return Event(iso8601.parse_date(e.time), e.type, json.loads(e.content))

    def next(self):
        req = backtests.BacktestEventRequest(id=self._id)
        self._next_queue.put(req)


class Backtest(object):
    def __init__(
        self,
        start_time: datetime,
        end_time: datetime = datetime.now(),
        accounts: Dict[str, Account] = {
            "binance": Account({"USDC": 100000}),
        },
    ):
        self._config = Config()
        self._candlesticks = Candlesticks()
        self._channel = grpc.insecure_channel(self._config.backtests_url)
        self._stub = backtests_grpc.BacktestsServiceStub(self._channel)
        self._id = self._create(start_time, end_time, accounts)
        self._start_time = start_time
        self._actual_time = self._start_time
        self._end_time = end_time

    def _create(self, start_time, end_time, accounts) -> int:
        if start_time.tzinfo is None or start_time.tzinfo.utcoffset(start_time) is None:
            raise Exception("no timezone specified on start")

        return self._stub.CreateBacktest(
            backtests.CreateBacktestRequest(
                start_time=start_time.isoformat(),
                end_time=end_time.isoformat(),
                accounts=self._account_to_grpc(accounts),
            )
        ).id

    def _account_to_grpc(self, accounts: Dict[str, Account]):
        req_accounts = {}
        for exch, account in accounts.items():
            assets = {}
            for asset, quantity in account.assets.items():
                assets[asset] = quantity
            req_accounts[exch] = backtests.Account(assets=assets)
        return req_accounts

    def subscribe(self, exchange_name, pair_symbol):
        self._stub.SubscribeToBacktestEvents(
            backtests.SubscribeToBacktestEventsRequest(
                id=self._id,
                exchange_name=exchange_name,
                pair_symbol=pair_symbol,
            )
        )

    def listen(self) -> BacktestEvents:
        return BacktestEvents(self._id, self._stub)

    def order(self, type: str, exchange: str, pair: str, side: str, quantity: float):
        req = backtests.CreateBacktestOrderRequest(
            backtest_id=self._id,
            type=type,
            exchange_name=exchange,
            pair_symbol=pair,
            side=side,
            quantity=quantity,
        )
        self._stub.CreateBacktestOrder(req)

    def accounts(self) -> Dict[str, Account]:
        req = backtests.AccountsRequest(
            backtest_id=self._id,
        )
        resp = self._stub.Accounts(req)
        return self._grpc_to_accounts(resp)

    def _grpc_to_accounts(self, resp: backtests.AccountsResponse) -> Dict[str, Account]:
        accounts = {}
        for exch, account in resp.accounts.items():
            assets = {}
            for asset, quantity in account.assets.items():
                assets[asset] = quantity
            accounts[exch] = Account(assets)
        return accounts

    def orders(self) -> List[Order]:
        req = backtests.OrdersRequest(backtest_id=self._id)
        resp = self._stub.Orders(req)
        return self._grpc_orders(resp)

    def _grpc_orders(self, resp: backtests.OrdersResponse) -> List[Order]:
        orders = []
        for o in resp.orders:
            orders.append(
                Order(
                    iso8601.parse_date(o.time),
                    o.type,
                    o.exchange_name,
                    o.pair_symbol,
                    o.side,
                    o.quantity,
                    o.price,
                )
            )
        return orders

    def on_event(self, event: Event):
        pass

    def on_end(self):
        pass

    def display(self, exchange: str, pair: str, period: Period):
        p = Grapher()

        start = self._start_time
        end = self._end_time
        cs = self._candlesticks.get(exchange, pair, period, start, end)
        p.candlesticks(cs)

        p.orders(self.orders())

        p.show()

    def actual_time(self) -> datetime:
        return self._actual_time

    def run(self):
        events = self.listen()
        finished = False
        while finished is False:
            events.next()

            while True:
                event = events.get()

                if event.type == "status":
                    finished = event.content["finished"]
                    self._actual_time = event.time
                    break

                self.on_event(event)

        return self.on_end()

    def candlesticks(
        self,
        exchange: str,
        pair: str,
        period: Period,
        relative_start: int,
        relative_end: int = 0,
        limit: int = 0,
    ):
        start = self._actual_time - relative_start * period.duration()
        end = self._actual_time - relative_end * period.duration()
        return self._candlesticks.get(exchange, pair, period, start, end, limit)
