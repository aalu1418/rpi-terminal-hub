from eufy import Eufy
from weather import Weather, clear
from webTime import WebTime
from time import time, sleep
import os
from datetime import datetime

class Loop():
    def __init__(self, location='Toronto', filename='eufy.json', increment=15, autorun=True, schedule={}):
        self.startup()
        self.webTime = WebTime()
        self.weather = Weather(location)
        self.eufy = Eufy(filename=filename)
        self.increment = increment
        self.schedule = {k:datetime.strptime(schedule[k],"%I%p").hour for k in schedule.keys()}
        self.runToday = False

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

        # check at midnight if vacuum is supposed to be run today
        if hour == 0 and minute < 5:
            if weekday in self.schedule.keys():
                self.runToday = True

        # if vacuum is supposed to be run today, check for the correct time
        if self.runToday == True and hour == self.schedule[weekday] and minute < 5:
            self.runToday = False
            return True

        return False

    def loop(self):
        clear() #clear loaded data
        while True:
            self.weather.fetch()
            self.webTime.fetch()

            # output info to CLI
            print(self.weather.data)
            print(f"Last updated: {self.webTime.timestamp}")

            # scheduled tasks
            if scheduler():
                self.eufy.emit('start_stop')
                print("Eufy Started")

            # pause until next interval
            self.delayCalc()
            sleep(self.delay)
            clear()



if __name__ == '__main__':
    schedule = {"Tuesday": "6PM",
                "Thursday": "6PM",
                "Saturday": "1PM"}
    loop = Loop(schedule=schedule)
