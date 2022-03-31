class Account(object):

    def __init__(self, balances: dict = {}):
        self.balances = balances

    def __repr__(self):
        return str(self.__dict__)
