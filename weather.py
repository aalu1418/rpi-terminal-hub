import requests, time

def clear():
    print("\033c", end="")

class Weather():
    def __init__(self, location):
        super(Weather, self).__init__()
        self.location = location

    def fetch(self):
        res = requests.get(f"https://wttr.in/{self.location}")
        self.data = res.text
        self.timestamp = time.strftime("%b %d @ %I:%M %p")

if __name__ == '__main__':
    weather = Weather('Toronto')
    clear()
    while True:
        weather.fetch()
        print(weather.data, weather.timestamp)
        time.sleep(15*60)
        clear()
