class Account(object):

    def __init__(self, assets: dict = {}):
        self.assets = assets

    def __repr__(self):
        return str(self.__dict__)
