# accomodates different usage scopes for modules
try:
    from modules import *  # scope for running alerts.py directly
except Exception as e:
    from src.modules import *  # scope for running alerts in main.py

options = {
    "toronto, on": ttc.TTC,
    "null, null": null.NULL,
}


class Alerts:
    def __init__(self, location):
        # default to NOAA
        if location.lower() not in options:
            options[location.lower()] = noaa.NOAA

        self.submodule = options[location.lower()]()
        self.name = self.submodule.name

    def fetch(self):
        # fetch in submodule and pass back data
        self.submodule.fetch()
        self.data = self.submodule.data


if __name__ == "__main__":
    alerts = Alerts("toronto, on")
    alerts.fetch()
    print(alerts.data)
