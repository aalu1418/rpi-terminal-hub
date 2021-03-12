# example: https://fishandwhistle.net/post/2016/raspberry-pi-pure-python-infrared-remote-control/

from gpiozero import LineSensor
from signal import pause
import time

sensor = LineSensor(17, sample_rate=38000, queue_len=1)
# sensor.when_line = lambda: print('Line detected')
# sensor.when_no_line = lambda: print('No line detected')
# pause()

data = []
sensor.wait_for_line()

start = time.time()
for i in range(15000):
    data.append(sensor.value)

print(sum(data), len(data), time.time()-start)
print("".join([str(i) for i in data]))



# if __name__ == "__main__":
