from datetime import datetime


class Tick(object):

    def __init__(self, time: datetime, exchange: str, pair: str, price: float):
        self.time = time
        self.exchange = exchange
        self.pair = pair
        self.price = price