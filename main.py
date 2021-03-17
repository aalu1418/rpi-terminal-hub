# custom modules
from webTime import WebTime
from weather import Weather, clear
# from eufy import Eufy
from ttc import TTC

# public modules
from time import time, sleep
import os
from datetime import datetime
import sys

# flask server
from flask import Flask, render_template
app = Flask(__name__)
app.output_data = "Hello, World" # storing data for output

# ------ MAIN LOOP ------------
class Loop():
    def __init__(self, schedule={}, cli=True, output=None, location='Toronto', filename='/home/pi/rpi-terminal-hub/eufy.json', increment=15, autorun=True):
        # self.startup()
        self.webTime = WebTime()
        self.weather = Weather(location, server=not cli)
        # self.eufy = Eufy(filename=filename)
        self.ttc = TTC()

        self.increment = increment
        self.schedule = {k:datetime.strptime(schedule[k],"%I%p").hour for k in schedule.keys()}
        self.runToday = False
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
            if weekday in self.schedule.keys():
                self.runToday = True

        # if vacuum is supposed to be run today, check for the correct time
        if self.runToday == True:
            if hour == self.schedule[weekday] and minute < 5:
                self.runToday = False
                return True

        return False

    def loop(self):
        clear() #clear loaded data
        while True:
            output = []

            try:
                self.weather.fetch()
                self.webTime.fetch()
                self.ttc.fetch()
            except Exception as e:
                print(e)
                return #retry fetch

            output.append(self.weather.data)
            output.append(f"Last updated: {self.webTime.timestamp}")
            output.append(", ".join(self.ttc.data))

            # scheduled tasks
            if self.scheduler():
                self.eufy.emit('start_stop')
                output.append("Eufy Started")

            if self.cli == False:
                data = output[0]
                data["updated"] = output[1]
                data["ttc"] = output[2] or None

                try:
                    data["eufy"] = output[3]
                    data["eufy"] = "Yes"
                except Exception as e:
                    data["eufy"] = "No"

                self.output.put(data)
            else:
                output = "\n".join(output)
                print(output)

            # clear start up triggers
            if self.start:
                self.start = False

            # pause until next interval
            self.delayCalc()
            sleep(self.delay)
            clear()

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
