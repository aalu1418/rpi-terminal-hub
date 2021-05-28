import json

class State:
    def __init__(self):
        self.changed = False

    def update(self, key, data):
        if not hasattr(self, key): # if new key, set
            setattr(self, key, data)
            self.changed = True

        self.changed = json.dumps(data) == json.dumps(getattr(self, key)) # check if existing value matchess
        if self.changed: # if check does not match, update
            setattr(self, key, data)

    def clear(self):
        self.changed = False
