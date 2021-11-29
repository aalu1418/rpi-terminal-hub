import requests, re, feedparser

# module for toronto transit commission
class NULL:
    def fetch(self):
        self.data = ["N/A"]


if __name__ == "__main__":
    ttc = NULL()
    ttc.fetch()
    print(ttc.data)
