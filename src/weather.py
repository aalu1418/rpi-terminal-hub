import requests, json, time
import re, os
from datetime import datetime
from dotenv import load_dotenv

load_dotenv()


def clear():
    print("\033c", end="")


iconMappings = {
    "01d": "sunny",
    "01n": "wi-night-clear",
    "02d": "sunny-overcast",
    "02n": "partly-cloudy",
    "03": "cloudy",
    "04": "wi-cloudy",
    "09": "showers",
    "10": "rain",
    "11": "thunderstorm",
    "13": "snow",
    "50d": "fog",
    "50n": "wi-night-fog",
}


def iconMapper(icon):
    prefix = "wi"
    dayNight = "day" if icon[-1] == "d" else "night-alt"

    try:
        name = iconMappings[icon[:2]]
    except Exception as e:  # if undefined try the whole name
        name = iconMappings[icon]

    # including the full name acts as an override for autogenerating the name
    return name if "wi-" in name else "-".join([prefix, dayNight, name])


class Weather:
    def __init__(self, locations, location=None):
        # use first location in json if key is not defined
        if location == None or location not in locations:
            location = list(locations.keys())[0]

        self.name = location
        self.location = locations[location]
        self.json = json

    def fetch(self):
        KEY = os.getenv("OWM_KEY")
        query = f"?lat={self.location['lat']}&lon={self.location['lon']}&appid={KEY}&units={self.location['units']}&exclude=minutely,daily"

        # units
        if self.location["units"] == "imperial":
            units = {"temp": "F", "precip": "in", "speed": "mph"}
        else:
            units = {"temp": "C", "precip": "mm", "speed": "km/h"}

        self.data = {"units": units}

        def parse(d):
            # handle missing precip data
            for i in ["rain", "snow"]:
                if i not in d.keys():
                    d[i] = {"1h": 0, "3h": 0}

            # handle missing wind data
            if "wind_gust" not in d.keys():
                d["wind_gust"] = 0

            # adjust wind units
            factor = 60 * 60 / 1000
            if self.location["units"] == "imperial":
                factor = 1

            date = datetime.fromtimestamp(d["dt"])
            wind = (
                f"{round(d['wind_speed'] * factor)}-{round(d['wind_gust'] * factor)}"
                if d["wind_gust"] != 0
                else round(d["wind_speed"] * factor)
            )

            output = {
                "condition": d["weather"][0]["description"].title(),
                "temp": round(d["temp"]),
                "temp_feel": d["feels_like"],
                "temp_feel_round": round(d["feels_like"]),
                "humidity": d["humidity"],
                "wind": wind,
                "precip": round(sum(d["rain"].values()) + sum(d["snow"].values()), 1),
                "dt": date.strftime("%a %d %b"),
                "hour": date.strftime("%-I %p"),
                "precip_percent": round(d["pop"] * 100),
                "iconPath": iconMapper(d["weather"][0]["icon"]),
                "uvi": d["uvi"],
                "clouds": d["clouds"],
            }

            # add sunrise/sunset time if present
            if "sunrise" in d.keys() and "sunset" in d.keys():
                output["sunrise"] = datetime.fromtimestamp(d["sunrise"]).strftime(
                    "%-I:%M %p"
                )
                output["sunset"] = datetime.fromtimestamp(d["sunset"]).strftime(
                    "%-I:%M %p"
                )

            return output

        # fetch onecall
        res = requests.get("https://api.openweathermap.org/data/2.5/onecall" + query)
        res = res.json()

        # current not always accurate, replace with first hourly (keeps params hourly is missing)
        res["current"].update(res["hourly"][0])

        # parse current
        self.data["current"] = parse(res["current"])
        self.data["location"] = self.name

        # parse forecast
        forecast = []
        for i in [res["hourly"][x] for x in range(1, 12, 3)]:
            forecast.append(parse(i))
        self.data["forecast"] = forecast


if __name__ == "__main__":
    from utils import getData

    locations = getData("locations.json")
    weather = Weather(locations)
    clear()
    while True:
        weather.fetch()
        print(weather.data)
        time.sleep(15 * 60)
        clear()
