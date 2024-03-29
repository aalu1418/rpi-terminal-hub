# accomodates different usage scopes for modules
try:
    from modules import *  # scope for running transit.py directly
except Exception as e:
    from src.modules import *  # scope for running transit in main.py

options = {"toronto, on": ttc.TTC}


class Transit:
    def __init__(self, location):
        # import module
        self.submodule = options[location.lower()]()

    def fetch(self):
        # fetch in submodule and pass back data
        self.submodule.fetch()
        self.data = self.submodule.data


if __name__ == "__main__":
    transit = Transit("toronto, on")
    transit.fetch()
    print(transit.data)
