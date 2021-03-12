# example: https://fishandwhistle.net/post/2016/raspberry-pi-pure-python-infrared-remote-control/

from gpiozero import LineSensor
from signal import pause

sensor = LineSensor(17, sample_rate=2000)
# sensor.when_line = lambda: print('Line detected')
# sensor.when_no_line = lambda: print('No line detected')
# pause()

data = []
sensor.wait_for_line()

for i in range(2000):
    data.append(sensor.value)

print(sum(data), len(data))



# if __name__ == "__main__":
