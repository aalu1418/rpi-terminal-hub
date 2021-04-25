# Create a CommandSet for your remote control
from ircodec.command import CommandSet
import sys
from time import sleep

states = {0: None,
          1: "Scheduled",
          2: "Started",
          3: "Completed"}

class Eufy:
    def __init__(self, filename=None, emitter=None, receiver=None):
        self.status = 0
        self.commands = ['start_stop', 'home', 'circle', 'edge', 'auto']
        if filename is not None:
            self.controller = CommandSet.load(filename)
        elif emitter is not None and receiver is not None:
            self.controller = CommandSet(name='Eufy', emitter_gpio=emitter, receiver_gpio=receiver)
        else:
            raise

    def pair(self):
        for c in self.commands:
            print('PAIR: ', c)
            self.controller.add(c)
            sleep(1)

        self.controller.save_as('eufy.json')

    def emit(self, v):
        if v == 'start_stop':
            self.status = 2
        self.controller.emit(v)

    def print(self):
        return states[self.status]


if __name__ == '__main__':
    if '--pair' in sys.argv:
        eufy = Eufy(emitter=27, receiver=17)
        eufy.pair()
    else:
        eufy = Eufy(filename='eufy.json')
        eufy.emit('start_stop')
