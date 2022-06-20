from datetime import datetime


class Event(object):
    def __init__(self, time: datetime, type: str, content: dict):
        self.time = time
        self.type = type
        self.content = content
