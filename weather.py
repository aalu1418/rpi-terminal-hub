import requests, json, time
from webTime import WebTime

def clear():
    print("\033c", end="")

class Weather():
    def __init__(self, location):
        self.location = location
        self.json = json

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

        #mission character replacements
        self.data = self.data.replace("\u2015", "\u2500") #fix horizontal lines
        self.data = self.data.replace("\u2019", "\u00B4") #fix horizontal lines

        #arrow replacement with directions
        unicode_arrows = [
            ("\u2192", "W"),
            ("\u2190", "E"),
            ("\u2191", "S"),
            ("\u2193", "N"),
            ("\u2196", "SE"),
            ("\u2197", "SW"),
            ("\u2198", "NW"),
            ("\u2199", "NE")]

        divider = "\u0020\u0020\n\u0020\u0020\u0020"
        for arrow in unicode_arrows:
            self.data = self.data.replace(arrow[0], arrow[1])
            #remove extra spaces
            if len(arrow[1]) == 2:
                split_header = self.data.split(divider)
                raw = split_header[1].split(arrow[1]) #split where directions were inserted
                new = [raw[0]]
                for e in raw[1:]:
                    raw_sub = e.split("\u2502") #split each element at vertical line
                    raw_sub[0] = raw_sub[0][:-1] #remove a space from the end of the first element
                    new = new + ["\u2502".join(raw_sub)] #rejoin
                split_header[1] = arrow[1].join(new) #rejoin
                self.data = divider.join(split_header)

if __name__ == '__main__':
    weather = Weather('Toronto', json=True)
    webTime = WebTime()
    clear()
    while True:
        weather.fetch()
        print(weather.data)
        webTime.fetch()
        print(f"Last updated: {webTime.timestamp}")
        time.sleep(15*60)
        clear()
