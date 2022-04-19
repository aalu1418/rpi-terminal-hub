import feedparser, re, requests


class XKCD:
    def __init__(self, path="xkcd.png"):
        self.url = "https://xkcd.com/rss.xml"
        self.path = path

    def fetch(self):
        data = feedparser.parse(self.url)
        regex = r'src="(\S+)"'
        url = re.findall(regex, data["entries"][0]["summary"])
        img_data = requests.get(url[0]).content
        with open(self.path, "wb") as handler:
            handler.write(img_data)


if __name__ == "__main__":
    xkcd = XKCD()
    xkcd.fetch()
