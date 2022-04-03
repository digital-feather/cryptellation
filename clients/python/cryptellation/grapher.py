import plotly.graph_objects as go
from plotly.graph_objs.scatter import Marker
import plotly.offline as py
import numpy as np
import pandas as pd
from plotly.offline import init_notebook_mode, iplot


class Grapher(object):

    def __init__(self):
        self._figure = None
        self._data = None
        self._shades = []

    def shade(self, start, stop, color: str = "black", opacity: float = 0.1):
        self._figure.add_vrect(
            x0=start,
            x1=stop,
            fillcolor=color,
            opacity=opacity,
            line_width=0,
        )

    def candlesticks(self, data):
        self._data = data
        chart_data = go.Candlestick(x=self._data.index,
                                    open=self._data['open'],
                                    high=self._data['high'],
                                    low=self._data['low'],
                                    close=self._data['close'])
        self._figure = go.Figure(data=[chart_data])

    def markers(self, series, symbol, color):
        min = self._data['low'].min()
        max = self._data['high'].min()

        data = pd.DataFrame(index=self._data.index)
        data['marker'] = np.nan
        data['symbol'] = 0
        data['color'] = np.nan

        for idx in self._data.index:
            if idx not in series:
                continue

            marker = self._data.loc[idx, 'high'] + (max - min)
            data.loc[idx] = {
                'marker': marker,
                'symbol': symbol,
                'color': color
            }

        trace = go.Scatter(x=data.index,
                           y=data['marker'],
                           mode='markers',
                           name='markers',
                           marker=Marker(size=20,
                                         symbol=data["symbol"],
                                         color=data["color"]))
        self._figure.add_trace(trace)

    def orders(self, orders):

        buy = pd.Series()
        sell = pd.Series()
        for order in orders:
            if order.side == "buy":
                buy.at[order.time] = True
            else:
                sell.at[order.time] = True

        self.markers(buy, 'triangle-up', 'green')
        self.markers(sell, 'triangle-down', 'red')

    def line(self, series, color: str = 'green'):
        trace = go.Scatter(x=self._data.index,
                           y=series,
                           line=dict(color=color, width=1))
        self._figure.add_trace(trace)

    def show(self):
        py.plot(self._figure)

    def save(self, path: str):
        layout = go.Layout(autosize=False, width=1920, height=1080)
        self._figure.update_layout(layout)
        self._figure.write_image(path)

    def jupyter_plot(self):
        init_notebook_mode(connected=True)
        iplot(self._figure)