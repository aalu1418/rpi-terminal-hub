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
* `fetchAPI()`: Pulls from OpenWeatherMap and returns struct for the server page. Requires a `.env` file with the contents:
```
OWM_KEY=<insert api key here>
```

### Eufy
`eufy.py`: Integration with IR emitter + receiver to control a Eufy Robovac 15C. The `eufy.json` file can be regenerated as well for other remotes.
* `pair()`: Sequence for recording various buttons on Eufy remote
* `emit()`: Command for emitting specified IR commands. Inputs should be a matching string to the stored commands.

### TTC
`ttc.py`: Checks TTC (toronto transit system) website for alerts using Regex - only displays alerts for subway and streetcars.
* `fetch()`: Pull HTML from [TTC alerts site](https://www.ttc.ca/Service_Advisories/all_service_alerts.jsp)
* `parse()`: Filter the HTML data using Regex into an array of TTC lines with active alerts

## Web Server
Endpoint: `GET /`
* Main endpoint for webpage with all information

Endpoint: `POST /vacuum`
* Example curl request: `curl -X POST -d "cmd=start" localhost:5000/vacuum`
* Missing `data` will default to `start_stop`
* Potential commands: `start_stop`, `home`, `circle`, `edge`, `auto`

Endpoint: `POST /pull`
* Runs `git pull` on the repository on remote server
* Example curl request: `curl -X POST localhost:5000/pull`

Endpoint `GET /logs`
* Returns the `cron.log` file
* Viewable using http://localhost:5000/logs or `curl -X POST -d "format=cli" localhost:5000/logs`
* Format can be `web` (default) or `cli`

Endpoint `POST /reboot`
* Restarts the server
* Example curl request: `curl -X POST -d "discardLog=True" localhost:5000/reboot`

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

RPi Timezones
* https://raspberrypi.stackexchange.com/questions/87164/setting-timezone-non-interactively

Server Inspiration
* https://shkspr.mobi/blog/2020/02/turn-an-old-ereader-into-an-information-screen-nook-str/
* Completed Look: ![](./media/NST_display.jpg)

To Do:
* Bash script for installing everything
* Post endpoint for triggering vacuum, and pulling new repo changes
