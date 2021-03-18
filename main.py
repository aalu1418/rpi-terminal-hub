# custom modules
from webTime import WebTime
from weather import Weather, clear
from eufy import Eufy
from ttc import TTC

# public modules
from time import time, sleep
import os
from datetime import datetime
import sys, traceback

# flask server
from flask import Flask, render_template
app = Flask(__name__)
app.output_data = "Hello, World" # storing data for output

# ------ MAIN LOOP ------------
class Loop():
    def __init__(self, schedule={}, cli=True, output=None, location='Toronto', filename='/home/pi/rpi-terminal-hub/eufy.json', increment=15, autorun=True):
        self.startup()
        self.webTime = WebTime()
        self.weather = Weather(location, server=not cli)
        self.eufy = Eufy(filename=filename)
        self.ttc = TTC()

        self.increment = increment
        self.schedule = {k:datetime.strptime(schedule[k],"%I%p").hour for k in schedule.keys()}
        self.cli = cli
        self.output = output
        self.start = True

        if autorun:
            self.loop()

    # initialization
    def startup(self):
        os.system("sudo pigpiod")

    # calculate seconds until next minute interval (aligned with the hour)
    def delayCalc(self):
        inc_sec = 60*self.increment
        time_sec = round(time())
        self.delay = inc_sec - time_sec % inc_sec

    # track schedule and see if task needs to be run
    def scheduler(self):
        weekday = self.webTime.weekday
        hour = self.webTime.raw.hour
        minute = self.webTime.raw.minute

        # check at midnight if vacuum is supposed to be run today (or on startup)
        if (hour == 0 and minute < 5) or self.start:
            self.eufy.status = 0 # reset status
            if weekday in self.schedule.keys():
                self.eufy.status = 1

        # if vacuum is supposed to be run today, check for the correct time
        if self.eufy.status == 1:
            if hour == self.schedule[weekday] and minute < 5:
                self.eufy.status = 2
                self.time = time()
                return True

        # if vacuum is running, after one hour mark as complete
        if self.eufy.status == 2 and time()-self.time > 60*60:
            self.eufy.status = 3

        return False

    def loop(self):
        while True:
            output = []

            try:
                self.weather.fetch()
                self.webTime.fetch()
                self.ttc.fetch()
            except Exception as e:
                traceback.print_exc()
                continue #retry fetch

            output.append(self.weather.data)
            output.append(f"Last Updated: {self.webTime.timestamp}")
            output.append(f"TTC Alerts: {', '.join(self.ttc.data) or None}")

            # scheduled tasks
            if self.scheduler():
                self.eufy.emit('start_stop')
            output.append(f"Eufy Status: {self.eufy.print()}")

            if self.cli == False:
                data = output[0]
                data["updated"] = output[1]
                data["ttc"] = output[2] or None
                data["eufy"] = output[3]
                self.output.put(data)
            else:
                clear() #clear loaded data
                output = "\n".join(output)
                print(output)

            # clear start up triggers
            if self.start:
                self.start = False

            # pause until next interval
            self.delayCalc()
            sleep(self.delay)

if __name__ == '__main__':
    cli = True
    schedule = {"Tuesday": "6PM",
                "Thursday": "6PM",
                "Saturday": "1PM"}

    # requires multiprocessing to run Flask server + loop
    if "--server" in sys.argv:
        from multiprocessing import Process, Queue
        output = Queue()
        cli = False
        p = Process(target=Loop, args=(schedule, cli, output))
        p.start()

        @app.route('/')
        def index():
            if not output.empty():
                app.output_data = output.get()

            # return initial string if data is not loaded
            if type(app.output_data) is str:
                return app.output_data

            # return app.output_data
            return render_template('index.html', data=app.output_data)

        app.run(host='0.0.0.0')
        p.join()

    else:
        loop = Loop(schedule)
