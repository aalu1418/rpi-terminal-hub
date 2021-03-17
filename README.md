# Raspberry Pi Terminal Hub
Turning a RPi into a hub using Raspian OS Lite (CLI only).

The `main.py` file contains the primary logic of continuously running the code. It handles updating the weather data (every 15 minutes), and also triggering scheduled events (specifically for the robot vacuum). It outputs data as a CLI display.
* `--server`: Outputs data as a web page on the local network. Implements the main loop and Flask webpage using the multiprocessing library (`Process` for running the loop, `Queue` for passing information from loop to web server).

## Modules
Various integrations for controlling / reporting devices and information.

### Web Time
`webTime.py`: Fetches time based on a web API and IP address location. Eliminates the need for calculating daylight savings, etc. It defaults to EST - timezones can be specified in the initialization (`timezone [string]`).
* `fetch()`: Pulls from time API for the given timezone

### Weather
`weather.py`: Fetches weather based on submitted location, and returns it via command line. Characters are corrected for proper display on RPi CLI. Takes a `location [string]` and `server [bool]` parameter
* `fetch()`: Pulls from [wttr.in](wttr.in) and formats the data string into a compatible CLI for Raspberry Pi
* `parse()`: Converts the CLI output to a struct/dictionary using Regex and data sorting.

### Eufy
`eufy.py`: Integration with IR emitter + receiver to control a Eufy Robovac 15C. The `eufy.json` file can be regenerated as well for other remotes.
* `pair()`: sequence for recording various buttons on Eufy remote
* `emit()`: command for emitting specified IR commands. Inputs should be a matching string to the stored commands.

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

Server Inspiration
* https://shkspr.mobi/blog/2020/02/turn-an-old-ereader-into-an-information-screen-nook-str/

ToDo:
* Need to fix API request robustness issues
