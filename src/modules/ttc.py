import requests, re, feedparser

# module for toronto transit commission
class TTC:
    def __init__(self):
        self.url = "http://ttc.ca/RSS/Service_Alerts/index.rss"
        self.name = "TTC"

    def fetch(self):
        self.data = feedparser.parse(self.url)
        self.parse()

    def parse(self):
        self.data = [
            alert["summary"].replace(":", "") for alert in self.data["entries"]
        ]

        # filter for streetcars + subways
        temp = [
            v[:6] for v in self.data if "Line" in v and "Elevator" not in v
        ]  # only pulls `Line #`
        temp = temp + [
            " ".join(v.split(" ")[0:2])
            for v in self.data
            if re.search(r"5[0-9]{2}[a-zA-Z]{0,1}.*", v)
        ]  # only pulls streetcars
        self.data = list(set(temp))


if __name__ == "__main__":
    ttc = TTC()
    ttc.fetch()
    print(ttc.data)
