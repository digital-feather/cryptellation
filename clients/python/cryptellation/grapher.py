import plotly.graph_objects as go
import plotly.offline as py
from datetime import datetime
import numpy as np
from ta.trend import SMAIndicator

from cryptellation.models.period import Period
from cryptellation.services.candlesticks import Candlesticks


class Grapher(object):

    def __init__(self):
        self._candlesticks = Candlesticks()
        self._figure = None
        self._data = None

    def candlesticks(self, exchange: str, pair: str, period: Period,
                     start: datetime, end: datetime):
        self._data = self._candlesticks.get(exchange, pair, period, start, end)
        chart_data = go.Candlestick(x=self._data.index,
                                    open=self._data['open'],
                                    high=self._data['high'],
                                    low=self._data['low'],
                                    close=self._data['close'])
        self._figure = go.Figure(data=[chart_data])

    def orders(self, orders):
        self._data['marker'] = np.nan
        self._data['symbol'] = 0
        self._data['color'] = np.nan

        min = self._data['low'].min()
        max = self._data['high'].min()

        for order in orders:
            self._data.at[order.time,
                          'marker'] = self._data.loc[order.time,
                                                     'high'] + (max - min)
            self._data.at[
                order.time,
                'symbol'] = "triangle-up" if order.side == "buy" else "triangle-down"
            self._data.at[order.time,
                          'color'] = "green" if order.side == "buy" else "red"

        trace = go.Scatter(x=self._data.index,
                           y=self._data['marker'],
                           mode='markers',
                           name='markers',
                           marker=go.Marker(size=20,
                                            symbol=self._data["symbol"],
                                            color=self._data["color"]))
        self._figure.add_trace(trace)

    def simple_moving_average(self, window: int, color: str = 'green'):
        sma = SMAIndicator(close=self._data['close'], window=window)
        data = sma.sma_indicator()

        trace = go.Scatter(x=self._data.index,
                           y=data,
                           line=dict(color=color, width=1))
        self._figure.add_trace(trace)

    def show(self):
        py.plot(self._figure)

    def save(self, path: str):
        layout = go.Layout(autosize=False, width=1920, height=1080)
        self._figure.update_layout(layout)
        self._figure.write_image(path)