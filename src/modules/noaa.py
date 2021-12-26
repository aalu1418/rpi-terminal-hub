import requests, re, feedparser

# module for NOAA
class NOAA:
    def __init__(self):
        # self.url = "https://alerts.weather.gov/cap/wwaatmget.php?x=COC069&y=0"
        self.url = "https://alerts.weather.gov/cap/wwaatmget.php?x=MIC161&y=0"
        self.name = "NOAA"

    def fetch(self):
        self.data = feedparser.parse(self.url)
        self.parse()

    def parse(self):
        self.data = [
            alert["title"].split(" issued ")[0]
            for alert in self.data["entries"]
            if "no active watches, warnings or advisories" not in alert["title"]
        ]

        # remove any duplicates
        self.data = list(set(self.data))


if __name__ == "__main__":
    noaa = NOAA()
    noaa.fetch()
    print(noaa.data)
