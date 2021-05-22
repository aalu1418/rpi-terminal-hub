import requests, json, time
import re, os
from datetime import datetime
from dotenv import load_dotenv
load_dotenv()

def clear():
    print("\033c", end="")

iconMappings = {
    '01d': 'sunny',
    '01n': 'wi-night-clear',
    '02d': 'sunny-overcast',
    '02n': 'partly-cloudy',
    '03': 'cloudy',
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
    def __init__(self, location, units='metric'):
        self.location = location
        self.json = json
        self.units = units

    def fetch(self):
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
    weather = Weather('Toronto')
    clear()
    while True:
        weather.fetch()
        print(weather.data)
        time.sleep(15*60)
        clear()
