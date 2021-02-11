import requests, time, json

def clear():
    print("\033c", end="")

class Weather():
    def __init__(self, location):
        super(Weather, self).__init__()
        self.location = location

    def fetch(self):
        res = requests.get(f"https://wttr.in/{self.location}")
        # print(json.dumps(res.text))
        self.data = res.text.replace("\n\nFollow \033[46m\033[30m@igor_chubin\033[0m for wttr.in updates\n", "")

        #color replacement
        white = "\u001b[0;0;0;1m"
        self.data = self.data.replace("\u001b[38;5;240;1m", white) #replace dark grey clouds with white
        self.data = self.data.replace("\u001b[38;5;250m", white) #replace grey clouds with white
        self.data = white + self.data #show white text
        self.data = self.data.replace("\u001b[0m", "\u001b[0m"+white) #show white text

        #arrow replacement with directions
        self.timestamp = time.strftime("%b %d @ %I:%M %p")

if __name__ == '__main__':
    weather = Weather('Toronto')
    clear()
    while True:
        weather.fetch()
        print(weather.data)
        print(f"Last updated: {weather.timestamp}")
        time.sleep(15*60)
        clear()
