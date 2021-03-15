# Raspberry Pi Terminal Hub
Turning a RPi into a hub using Raspian OS Lite (CLI only).

The `main.py` file contains the primary logic of continuously running the code. It handles updating the weather data (every 15 minutes), and also triggering scheduled events (specifically for the robot vacuum).

## Modules
Various integrations for controlling / reporting devices and information.

### Web Time
`webTime.py`: Fetches time based on a web API and IP address location. Eliminates the need for calculating daylight savings, etc.

### Weather
`weathery.py`: Fetches weather based on submitted location, and returns it via command line. Characters are corrected for proper display on RPi CLI.

### Eufy
`eufy.py`: Integration with IR emitter + receiver to control a Eufy Robovac 15C. The `eufy.json` file can be regenerated as well for other remotes.

## Notes
Autorun on RPi using CRON
* https://phoenixnap.com/kb/crontab-reboot
* Included a delay to allow internet connection to be established
```
@reboot sleep 60 && /usr/bin/python3 /home/pi/rpi-terminal-hub/main.py --server >> ~/cron.log 2>&1
```

Terminal Weather
* https://github.com/chubin/wttr.in

IR Receiver + Emitter
* (Circult) https://www.hackster.io/austin-stanton/creating-a-raspberry-pi-universal-remote-with-lirc-2fd581
* (Python Package) https://github.com/kentwait/ircodec
