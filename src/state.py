import json

class State:
    def __init__(self):
        self.changed = False

    def update(self, key, data):
        if not hasattr(self, key): # if new key, set
            setattr(self, key, data)
            self.changed = True
            return

        self.changed = json.dumps(data) != json.dumps(getattr(self, key)) # check if existing value matchess
        if self.changed: # if check does not match, update
            setattr(self, key, data)

    def clear(self):
        self.changed = False

if __name__ == '__main__':
    def printVar(func=None, *args):
        if func is not None:
            func(*args)
        print('Update:', test.changed, test.__dict__)

    test = State()
    params = [[None],
    [test.update, 'something', 1],
    [test.clear],
    [test.update, 'something', 1],
    [test.update, 'something', 2]]

    for param in params:
        printVar(*param)
