# example: https://fishandwhistle.net/post/2016/raspberry-pi-pure-python-infrared-remote-control/

from gpiozero import PWMLED
import json
import time

class Emitter:
    def __init__(self, GPIO, filename=None):
        self.led = PWMLED(GPIO, frequency=38000)
        self.data = None

        if filename is not None:
            try:
                file = open(filename, 'r')
                raw = file.read()
                file.close()
                self.data = json.loads(raw)
            except Exception as e:
                print("ERROR LOADING FILE:", e)

    def run(self):
        if self.data is None:
            self.blink()
        else:
            self.emit()

    # blink based on file
    def emit(self):
        for i in range(len(self.data["0"])):
            self.led.on()
            time.sleep(self.data["0"][i])
            self.led.off()
            time.sleep(self.data["1"][i])

    # blink for 10 seconds
    def blink(self):
        for i in range(5):
            self.led.on()
            time.sleep(1)
            self.led.off()
            time.sleep(1)

if __name__ == '__main__':
    emitter = Emitter(18, "./data/IR_receive.json")
    emitter.run()
