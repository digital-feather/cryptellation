from datetime import datetime


class Event(object):

    def __init__(self, time: datetime, type: str, content: dict):
        self._time = time
        self._type = type
        self._content = content

    def time(self) -> datetime:
        return self._time

    def type(self) -> str:
        return self._type

    def content(self) -> dict:
        return self._content