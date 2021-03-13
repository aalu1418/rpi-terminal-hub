# Create a CommandSet for your remote control
from ircodec.command import CommandSet
import sys
from time import sleep

class Eufy:
    def __init__(self, filename=None, emitter=None, receiver=None):
        if filename is not None:
            self.controller = CommandSet.load(filename)
        elif emitter is not None and receiver is not None:
            self.controller = CommandSet(name='Eufy', emitter_gpio=emitter, receiver_gpio=receiver)
        else:
            raise

    def pair(self):
        commands = ['start_stop', 'home', 'circle', 'edge', 'auto']
        for c in commands:
            print('PAIR: ', c)
            self.controller.add(c)
            sleep(1)

        self.controller.save_as('eufy.json')

    def emit(self, v):
        self.controller.emit(v)

if __name__ == '__main__':
    if '--pair' in sys.argv:
        eufy = Eufy(emitter=18, receiver=17)
        eufy.pair()
    else:
        eufy = Eufy(filename='eufy.json')
        eufy.emit('start_stop')
