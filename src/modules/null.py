import requests, re, feedparser

# null module for alerts
class NULL:
    def __init__(self):
        self.name = ""

    def fetch(self):
        self.data = ["N/A"]


if __name__ == "__main__":
    ttc = NULL()
    ttc.fetch()
    print(ttc.data)
