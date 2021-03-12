# example: https://fishandwhistle.net/post/2016/raspberry-pi-pure-python-infrared-remote-control/

from gpiozero import LineSensor
from signal import pause

sensor = LineSensor(17)
sensor.when_line = lambda: print('Line detected')
sensor.when_no_line = lambda: print('No line detected')
pause()



# if __name__ == "__main__":
