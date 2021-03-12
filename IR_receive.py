# example: https://fishandwhistle.net/post/2016/raspberry-pi-pure-python-infrared-remote-control/

from gpiozero import LineSensor
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

duration = time.time()-start
print("Total duration: ", duration, "\n")

value = 0
parsed = {0: [], 1:[]}

# sorting data into struct
count = 0
for v in data:
    if v == value:
        count += 1
    else:
        parsed[value].append(count*duration/len(data))
        count = 0
        value = int(not(value))

# add last count to data struct
parsed[value].append(count*duration/len(data))

for i in range(len(parsed[0])):
    print("On: ", parsed[0][i])
    print("Off: ", parsed[1][i], "\n")
