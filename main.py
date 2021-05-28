# custom modules
from src.webTime import WebTime
from src.weather import Weather, clear
from src.transit import Transit
# from src.eufy import Eufy # imported dynamically in the Loop() and '__main__'

# public modules
from time import time, sleep
from datetime import datetime
import os, sys, traceback, subprocess
from multiprocessing import Process, Queue

# flask server
from flask import Flask, render_template, request
app = Flask(__name__)
app.output_data = "Hello, World" # storing data for output

# ------ MAIN LOOP ------------
class Loop():
    def __init__(self, schedule={}, output=None, eufy=True, location='Toronto', filename='/home/pi/rpi-terminal-hub/data/eufy.json', increment=15, autorun=True):
        if eufy:
            from src.eufy import Eufy
            self.startup()
            self.eufy = Eufy(filename=filename)
        self.webTime = WebTime()
        self.weather = Weather(location)
        self.transit = Transit(location)

        self.increment = increment
        self.schedule = {k:datetime.strptime(schedule[k],"%I:%M%p") for k in schedule.keys()}
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
        try:
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
                if hour == self.schedule[weekday].hour and 0 <= self.schedule[weekday].minute-minute < 15:
                    self.eufy.status = 2
                    self.time = time()
                    return True

            # if vacuum is running, after one hour mark as complete
            if self.eufy.status == 2 and time()-self.time > 60*60:
                self.eufy.status = 3
        except Exception as e:
            self.eufy.status = 0 # reset status on error & log
            traceback.print_exc()

        return False

    def loop(self):
        while True:
            output = []

            try:
                self.weather.fetch()
                self.webTime.fetch()
                self.transit.fetch()
            except Exception as e:
                traceback.print_exc()
                sleep(30) #pause 30 seconds
                continue #retry fetch

            output.append(self.weather.data)
            output.append(f"Last Updated: {self.webTime.timestamp}")
            output.append(f"Transit/Traffic Alerts: {', '.join(self.transit.data) or None}")

            # scheduled tasks
            if hasattr(self, 'eufy') and self.scheduler():
                self.eufy.emit('start_stop')
            output.append(f"Eufy Status: {self.eufy.print() if hasattr(self, 'eufy') else 'N/A'}")

            data = output[0]
            data["updated"] = output[1]
            data["transit"] = output[2] or None
            data["eufy"] = output[3]
            self.output.put(data)

            # clear start up triggers
            if self.start:
                self.start = False

            # pause until next interval
            self.delayCalc()
            sleep(self.delay)

if __name__ == '__main__':
    schedule = {"Tuesday": "6:00PM",
                "Thursday": "7:15PM",
                "Saturday": "1:00PM"}

    eufy = False
    if "--no-eufy" not in sys.argv:
        from src.eufy import Eufy
        eufy = True

    # requires multiprocessing to run Flask server + loop
    output = Queue()
    p = Process(target=Loop, args=(schedule, output, eufy))
    p.start()

    # pull data from queue until latest
    def pullLatest():
        while not output.empty():
            app.output_data = output.get()

    @app.route('/')
    def index():
        pullLatest()
        # return initial string if data is not loaded
        if type(app.output_data) is str:
            return app.output_data

        # return app.output_data
        return render_template('index.html', data=app.output_data)

    @app.route('/raw', methods=['GET'])
    def returnRaw():
        pullLatest()
        return app.output_data

    #  allow remote triggering of vacuum
    @app.route('/vacuum', methods=['POST'])
    def vacuum():
        if eufy:
            eufyInterface = Eufy(filename='/home/pi/rpi-terminal-hub/data/eufy.json')
            cmd = request.form.get('cmd').lower()
            if (cmd == 'start' or cmd == 'stop' or cmd not in eufyInterface.commands):
                cmd = 'start_stop'

            eufyInterface.emit(cmd)
            return {'status': 'called '+cmd}
        return {'status': 'eufy not enabled'}

    # allow remote pull to update code
    @app.route('/pull', methods=['POST'])
    def pull():
       out = os.system('cd /home/pi/rpi-terminal-hub && git pull')
       return {'status': 'ready for reboot'}

    # allow remote pull to update code
    @app.route('/logs', methods=['GET'])
    def logs():
        f = open('/home/pi/cron.log')
        data = f.read()
        f.close()

        format = request.form.get('format')
        if format == None or format.lower() == 'web':
            data = data.replace('\n', '<br>')

        return data

    # allow reboot
    @app.route('/reboot', methods=['POST'])
    def reboot():
        discardLog = request.form.get('discardLog')
        if discardLog == "True":
            subprocess.Popen('sleep 5 && rm /home/pi/cron.log', shell=True)

        subprocess.Popen('sleep 10 && sudo reboot', shell=True)
        return {'status': 'rebooting in 10 seconds'}


    app.run(host='0.0.0.0')
    p.join()
