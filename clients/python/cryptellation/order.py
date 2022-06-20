from datetime import datetime


class Order(object):
    def __init__(
        self,
        time: datetime,
        type: str,
        exchange: str,
        pair: str,
        side: str,
        quantity: float,
        price: float,
    ):
        self.time = time
        self.type = type
        self.exchange = exchange
        self.pair = pair
        self.side = side
        self.quantity = quantity
        self.price = price
