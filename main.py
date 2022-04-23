# custom modules
from src.webTime import WebTime
from src.weather import Weather, clear
from src.alerts import Alerts
from src.state import State
from src.utils import Cmd, getData

# from src.eufy import Eufy # imported dynamically in the Loop() and '__main__'

# public modules
from time import time, sleep
from datetime import datetime
import os, sys, traceback, subprocess
from multiprocessing import Process, Queue

# flask server
from flask import Flask, render_template, request

app = Flask(__name__)
app.output_data = "Server is starting... "  # storing data for output
wifi_down = "ERROR: failed to fetch data [wifi may be down]"

# ------ MAIN LOOP ------------
class Loop:
    def __init__(
        self,
        output=None,
        eufy=True,
        location=None,
        filename="/home/pi/rpi-terminal-hub/data/eufy.json",
        increment=15,
        autorun=True,
    ):
        schedule = getData("schedule.json")
        locations = getData("locations.json")

        if eufy:
            from src.eufy import Eufy

            self.startup()
            self.eufy = Eufy(filename=filename)
        self.webTime = WebTime()
        self.weather = Weather(locations, location)
        self.alerts = Alerts(location or list(locations.keys())[0])

        self.retry = False
        self.increment = increment
        self.schedule = {
            k: datetime.strptime(schedule[k], "%I:%M%p") for k in schedule.keys()
        }
        self.output = output
        self.start = True

        self.data = State()

        if autorun:
            self.loop()

    # initialization
    def startup(self):
        os.system("sudo pigpiod")

    # calculate seconds until next minute interval
    def delayCalc(self):
        time_sec = round(time())
        self.delay = 60 - time_sec % 60

    # track schedule and see if task needs to be run
    def scheduler(self):
        try:
            weekday = self.webTime.weekday
            hour = self.webTime.hour
            minute = self.webTime.minute

            # reset status at midnight or on start
            if (hour == 0 and minute == 0) or self.start:
                self.eufy.status = 0  # reset status

            # if there a schedule for today and it is not scheduled, change to scheduled
            if weekday in self.schedule.keys() and self.eufy.status == 0:
                self.eufy.status = 1

            # if vacuum is supposed to be run today, check for the correct time
            if self.eufy.status == 1:
                if (
                    hour == self.schedule[weekday].hour
                    and minute == self.schedule[weekday].minute
                ):
                    self.eufy.status = 2
                    self.time = time()
                    return True

            # if vacuum is running, after 1.5 hours mark as complete
            if self.eufy.status == 2 and time() - self.time > 90 * 60:
                self.eufy.status = 3
        except Exception as e:
            self.eufy.status = 0  # reset status on error & log
            traceback.print_exc()

        return False

    # fetch data with retry logic
    def fetch(self):
        try:
            # fetch data
            self.weather.fetch()
            self.webTime.fetch()
            self.alerts.fetch()

            # write data to state
            self.data.update("weather", self.weather.data)
            self.data.update("updated", self.webTime.timestamp)
            self.data.update(
                "alerts",
                f"Alerts{' ['+self.alerts.name+']' if self.alerts.name != '' else self.alerts.name}: {', '.join(self.alerts.data) or None}",
            )

            self.retry = False
        except Exception as e:  # if error, retry in the next minute
            self.retry = True

    def loop(self):
        while True:
            # fetch data on first run or retry request or based on minute interval
            if self.start or self.retry or self.webTime.minute % self.increment == 0:
                self.fetch()

            # run vacuum checks
            if hasattr(self, "eufy") and self.scheduler():
                self.eufy.emit("30min")
            self.data.update(
                "eufy",
                f"Eufy Status: {self.eufy.print() if hasattr(self, 'eufy') else 'N/A'}",
            )

            # if data is changed, push data to queue
            if self.data.changed:
                self.output.put(self.data.__dict__)  # return data
                self.data.clear()  # clear changed state

            # clear start up triggers
            if self.start:
                self.start = False

            # pause until next interval
            self.delayCalc()
            self.webTime.inc()  # increment minute (not polling API every minute)
            sleep(self.delay)


if __name__ == "__main__":
    eufy = False
    if "--no-eufy" not in sys.argv:
        from src.eufy import Eufy

        eufy = True

    # requires multiprocessing to run Flask server + loop
    output = Queue()
    p = Process(target=Loop, args=(output, eufy))
    p.start()

    # pull data from queue until latest
    def pullLatest():
        while not output.empty():
            app.output_data = output.get()

    @app.route("/")
    def index():
        pullLatest()
        # return initial string if data is not loaded
        if type(app.output_data) is str:
            return app.output_data

        if "weather" not in app.output_data:
            print("ERROR: failed no data to return")
            return wifi_down

        # return app.output_data
        return render_template("index.html", data=app.output_data)

    @app.route("/raw", methods=["GET"])
    def returnRaw():
        pullLatest()
        return app.output_data

    #  allow remote triggering of vacuum
    @app.route("/vacuum", methods=["POST"])
    def vacuum():
        if eufy:
            eufyInterface = Eufy(filename="/home/pi/rpi-terminal-hub/data/eufy.json")
            cmd = request.form.get("cmd")
            if cmd == "start" or cmd == "stop" or cmd not in eufyInterface.commands:
                cmd = "start_stop"

            eufyInterface.emit(cmd)
            return {"status": "called " + cmd}
        return {"status": "eufy not enabled"}

    # allow remote pull to update code
    @app.route("/pull", methods=["POST"])
    def pull():
        c = Cmd()
        os.system(c.pull)
        return {"status": "ready for reboot"}

    # allow remote pull to update code
    @app.route("/logs", methods=["GET"])
    def logs():
        f = open("/home/pi/cron.log")
        data = f.read()
        f.close()

        format = request.form.get("format")
        if format == None or format.lower() == "web":
            data = data.replace("\n", "<br>")

        return data

    # allow reboot
    @app.route("/reboot", methods=["POST"])
    def reboot():
        c = Cmd()
        subprocess.Popen(c.rmLog, shell=True)
        subprocess.Popen(c.reboot, shell=True)
        return {"status": "rebooting in 10 seconds"}

    # called before pushing to have server auto pull + reboot
    @app.route("/pre-push", methods=["POST"])
    def prePush():
        c = Cmd(15)
        seq = [c.pull, c.rmLog, c.reboot]
        for s in seq:
            subprocess.Popen(s, shell=True)
        return {"status": "ready for push"}

    app.run(host="0.0.0.0")
    p.join()
