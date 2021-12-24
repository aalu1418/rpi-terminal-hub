# accomodates different usage scopes for modules
try:
    from modules import *  # scope for running transit.py directly
except Exception as e:
    from src.modules import *  # scope for running transit in main.py

options = {"toronto, on": ttc.TTC}


class Alerts:
    def __init__(self, location):
        # import module
        if location.lower() not in options:
            options[location.lower()] = null.NULL

        self.submodule = options[location.lower()]()

    def fetch(self):
        # fetch in submodule and pass back data
        self.submodule.fetch()
        self.data = self.submodule.data


if __name__ == "__main__":
    alerts = Alerts("toronto, on")
    alerts.fetch()
    print(alerts.data)
