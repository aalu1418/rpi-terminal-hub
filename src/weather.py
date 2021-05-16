import requests, json, time
import re, os
from datetime import datetime

def clear():
    print("\033c", end="")

iconMappings = {
    '01d': 'sunny',
    '01n': 'wi-night-clear',
    '02': 'cloudy',
    '03': 'wi-cloud',
    '04': 'wi-cloudy',
    '09': 'showers',
    '10': 'rain',
    '11': 'thunderstorm',
    '13': 'snow',
    '50d': 'fog',
    '50n': 'wi-night-fog'
}

def iconMapper(icon):
    prefix = 'wi'
    dayNight = 'day' if icon[-1] is 'd' else 'night-alt'

    try:
        name = iconMappings[icon[:2]]
    except Exception as e: # if undefined try the whole name
        name = iconMappings[icon]

    # including the full name acts as an override for autogenerating the name
    return name if 'wi-' in name else '-'.join([prefix, dayNight, name])

class Weather():
    def __init__(self, location, server=False, units='metric'):
        self.location = location
        self.json = json
        self.server = server
        self.units = units

    def fetch(self):
        if self.server == True:
            from dotenv import load_dotenv
            load_dotenv()
            return self.fetchAPI()

        res = requests.get(f"https://wttr.in/{self.location}")
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
                raw = split_header[-1].split(arrow[1]) #split where directions were inserted
                new = [raw[0]]
                for e in raw[1:]:
                    raw_sub = e.split("\u2502") #split each element at vertical line
                    raw_sub[0] = raw_sub[0][:-1] #remove a space from the end of the first element
                    new = new + ["\u2502".join(raw_sub)] #rejoin
                split_header[-1] = arrow[1].join(new) #rejoin
                self.data = divider.join(split_header)

    def fetchAPI(self):
        KEY = os.getenv('OWM_KEY')
        query=f"?q={self.location}&appid={KEY}&units={self.units}"

        # units
        if self.units == 'imperial':
            units = {'temp': 'F',
                     'precip': 'mm',
                     'speed': 'mph'}
        else:
             units = {'temp': 'C',
                      'precip': 'mm',
                      'speed': 'km/h'}

        self.data = {'units': units}

        def parse(d):
            # handle missing precip data
            for i in ['rain', 'snow']:
                if i not in d.keys():
                    d[i] = {'1h': 0, '3h': 0}

            # handle missing probability of precipitation in current weather
            if "pop" not in d.keys():
                d["pop"] = 0

            #handle missing wind data
            if "wind" in d.keys() and "gust" not in d["wind"].keys():
                d["wind"]["gust"] = 0

            # adjust wind units
            factor = 60*60/1000
            if self.units == 'imperial':
                factor = 1

            date = datetime.fromtimestamp(d["dt"])

            return {"condition": d["weather"][0]["description"].title(),
                    "temp": round(d["main"]["temp"]),
                    "temp_feel": d["main"]["feels_like"],
                    "temp_feel_round": round(d["main"]["feels_like"]),
                    "wind": round(d["wind"]["speed"]*factor),
                    "gust": round(d["wind"]["gust"]*factor),
                    "precip": round(sum(d['rain'].values()) + sum(d['snow'].values()), 1),
                    "dt": date.strftime('%a %d %b'),
                    "hour": date.strftime('%I:%M %p'),
                    "precip_percent": round(d["pop"]*100),
                    "iconPath": iconMapper(d['weather'][0]['icon'])
                    }


        # fetch current
        res = requests.get("https://api.openweathermap.org/data/2.5/weather"+query)
        res = res.json()

        self.data["current"] = parse(res)
        self.data["location"] = f"Location: {res['name']} [{res['coord']['lat']}, {res['coord']['lon']}]"

        # fetch forecast
        res = requests.get("https://api.openweathermap.org/data/2.5/forecast"+query)
        res = res.json()
        res = res["list"][0:4]
        forecast = []
        for i in res:
            forecast.append(parse(i))
        self.data["forecast"] = forecast


if __name__ == '__main__':
    weather = Weather('Toronto', server=True)
    clear()
    while True:
        weather.fetch()
        print(weather.data)
        time.sleep(15*60)
        clear()
