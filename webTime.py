from datetime import datetime
import requests

weekdays = {1: "Monday",
            2: "Tuesday",
            3: "Wednesday",
            4: "Thursday",
            5: "Friday",
            6: "Saturday",
            7: "Sunday"}

class WebTime:
    def __init__(self):
        self.url = "http://worldtimeapi.org/api/ip"
        self.fetch()

    def fetch(self):
        res = requests.get("http://worldtimeapi.org/api/ip")
        res = res.json()
        self.raw = datetime.strptime(res['datetime'], "%Y-%m-%dT%H:%M:%S.%f%z")
        self.timestamp = self.raw.strftime("%b %d @ %I:%M %p")
        self.weekday = weekdays[self.raw.isoweekday()]

if __name__ == '__main__':
    webTime = WebTime()
    print(webTime.timestamp, webTime.weekday)
