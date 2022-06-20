from enum import Enum

import datetime


class Period(Enum):
    M1 = "M1"
    M3 = "M3"
    M5 = "M5"
    M15 = "M15"
    M30 = "M30"
    H1 = "H1"
    H2 = "H2"
    H4 = "H4"
    H6 = "H6"
    H8 = "H8"
    H12 = "H12"
    D1 = "D1"
    D3 = "D3"
    W1 = "W1"

    def __str__(self):
        return self.value

    def duration(self):
        if self in DURATIONS:
            return DURATIONS[self]
        else:
            raise ValueError("unknown period", self)


if __name__ == "__main__":
    print(Period.M1)

DURATIONS = {
    Period.M1: datetime.timedelta(minutes=1),
    Period.M3: datetime.timedelta(minutes=3),
    Period.M5: datetime.timedelta(minutes=5),
    Period.M15: datetime.timedelta(minutes=15),
    Period.M30: datetime.timedelta(minutes=30),
    Period.H1: datetime.timedelta(hours=1),
    Period.H2: datetime.timedelta(hours=2),
    Period.H4: datetime.timedelta(hours=4),
    Period.H6: datetime.timedelta(hours=6),
    Period.H8: datetime.timedelta(hours=8),
    Period.H12: datetime.timedelta(hours=12),
    Period.D1: datetime.timedelta(days=1),
    Period.D3: datetime.timedelta(days=3),
    Period.W1: datetime.timedelta(days=7),
}
