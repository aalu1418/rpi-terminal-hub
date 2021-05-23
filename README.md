# Raspberry Pi Terminal Hub
Turning a RPi into a server hub using Raspian OS Lite.

Credits: Icons created by [@erikflowers](https://github.com/erikflowers/weather-iconss)

The `main.py` file contains the primary logic of continuously running the code. It handles updating the weather data (every 15 minutes), and also triggering scheduled events (specifically for the robot vacuum). It outputs data as a Flask app at http://0.0.0.0:5000.
* `--no-eufy`: Useful flag when testing locally where the GPIO libraries are not installed. Skips the eufy related portions in the code.

## Modules
Various integrations for controlling / reporting devices and information.

### Web Time
`webTime.py`: Fetches time based on a web API and IP address location. Eliminates the need for calculating daylight savings, etc. It defaults to EST - timezones can be specified in the initialization (`timezone [string]`).
* `fetch()`: Pulls from time API for the given timezone

### Weather
`weather.py`: Fetches weather based on submitted location, and returns it as a data struct. Takes a `location [string]` parameter
* `fetch()`: Pulls from OpenWeatherMap and returns struct for the server page. Requires a `.env` file with the contents:
```
OWM_KEY=<insert api key here>
```

### Eufy
`eufy.py`: Integration with IR emitter + receiver to control a Eufy Robovac 15C. The `eufy.json` file can be regenerated as well for other remotes.
* `pair()`: Sequence for recording various buttons on Eufy remote
* `emit()`: Command for emitting specified IR commands. Inputs should be a matching string to the stored commands.

### Transit/Traffic
`transit.py`: Acts as a wrapper for passed in city parameter for specific transit or traffic modules defined for each city.
* `fetch()`: Runs the respective fetch command to retrieve and parse data.
* `data`: The parsed data
* Each submodule must have the following defined interface within the class:
   * `init(self)`: no parameters can be passed in, initializes the data needed
   * `fetch(self)`: no parameters, pulls and processes data into a `[]string`
   * `data`: Class parameter where the `[]string` is stored
   * Other functions may be included to assist with retrieving, processing, etc
   * See the [Toronto Transit (TTC) module](./src/modules/ttc.py) as an example

## Web Server
Endpoint: `GET /`
* Main endpoint for stylized webpage with all information

Endpoint: `GET /raw`
* Raw data used in stylized webpage
* Viewable using http://localhost:5000/raw or `curl -X GET localhost:5000/raw`

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
* Used in a previous version (useful for CLI weather app)

IR Receiver + Emitter
* (Circult) https://www.hackster.io/austin-stanton/creating-a-raspberry-pi-universal-remote-with-lirc-2fd581
* (Python Package) https://github.com/kentwait/ircodec

RPi Timezones
* https://raspberrypi.stackexchange.com/questions/87164/setting-timezone-non-interactively

Server Inspiration
* https://shkspr.mobi/blog/2020/02/turn-an-old-ereader-into-an-information-screen-nook-str/
* Completed Look (interface is a bit outdated): ![](./media/NST_display.jpg)

Ideas
* Add `traffic` module for driving commutes
