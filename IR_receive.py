# example: https://fishandwhistle.net/post/2016/raspberry-pi-pure-python-infrared-remote-control/

from gpiozero import LineSensor
from time import time
import json
import sys

class Receiver():
    def __init__(self, GPIO):
        self.GPIO = GPIO
        self.data = []
        self.parsed = {0: [], 1:[]}

        self.sensor = LineSensor(self.GPIO, sample_rate=38000, queue_len=1)

    # record + parse
    def run(self, n=15000):
        self.sensor.wait_for_line()

        start = time()
        for i in range(n):
            self.data.append(self.sensor.value)

        self.duration = time()-start

        # sorting data into struct
        value = 0
        count = 0
        for v in self.data:
            if v == value:
                count += 1
            else:
                self.parsed[value].append(count*self.duration/len(self.data))
                count = 0
                value = int(not(value))

        # add last count to data struct
        self.parsed[value].append(count*self.duration/len(self.data))

    # export data
    def print(self, export=False):
        print("Total duration: ", self.duration, "\n")

        for i in range(len(self.parsed[0])):
            print("On: ", self.parsed[0][i])
            print("Off: ", self.parsed[1][i], "\n")

        if export is True:
            file = open("./data/IR_receive.json", "w")
            data = json.dumps(self.parsed)
            file.writelines(data)
            file.close()

if __name__ == '__main__':
    export = False
    if "--export" in sys.argv:
        export = True

    receiver = Receiver(17)
    receiver.run()
    receiver.print(export)
