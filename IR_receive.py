# example: https://fishandwhistle.net/post/2016/raspberry-pi-pure-python-infrared-remote-control/

from gpiozero import LineSensor
from signal import pause

sensor = LineSensor(17, sample_rate=2000)
# sensor.when_line = lambda: print('Line detected')
# sensor.when_no_line = lambda: print('No line detected')
# pause()

while True:
    print(sensor.value)



# if __name__ == "__main__":
