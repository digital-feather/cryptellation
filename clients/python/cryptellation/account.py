from typing import Dict


class Account(object):
    def __init__(self, assets: Dict[str, float] = {}):
        self.assets = assets
