from datetime import datetime
import requests


class WebTime:
    def __init__(self, timezone="est"):
        self.url = f"http://worldclockapi.com/api/json/{timezone}/now"

    def fetch(self):
        res = requests.get(self.url)
        res = res.json()
        self.raw = datetime.strptime(res["currentDateTime"], "%Y-%m-%dT%H:%M%z")
        self.timestamp = self.raw.strftime("%b %d @ %I:%M %p")
        self.weekday = res["dayOfTheWeek"]
        self.minute = self.raw.minute

    def inc(self):
        self.minute = (self.minute + 1) % 60


if __name__ == "__main__":
    webTime = WebTime()
    webTime.fetch()
    print(webTime.timestamp, webTime.weekdays)
