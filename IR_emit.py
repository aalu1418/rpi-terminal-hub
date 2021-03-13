# example: https://fishandwhistle.net/post/2016/raspberry-pi-pure-python-infrared-remote-control/

from gpiozero import PWMLED
import time

led = PWMLED(18)

while True:
    led.on()
    time.sleep(1)
    led.off()
    time.sleep(1)
