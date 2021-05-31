import os, json
from itertools import compress
from datetime import datetime

# same as src.utils
def getData(name):
    dirname = os.path.dirname(__file__)  # use filepath to utils file
    datapath = os.path.realpath(os.path.join(dirname, "../data"))
    f = open(f"{datapath}/{name}")
    data = json.load(f)
    f.close()
    return data


def check(expect, actual, errMsg):
    if expect != actual:
        raise Exception(errMsg)


def checkStrFloat(n):
    check(True, isinstance(n, str), f"'{n}' is not a string")
    float(n)


def verifyLocations(data):
    l = list(data.keys())
    lKeys = ["lon", "lat", "units"]
    units = ["standard", "metric", "imperial"]

    # validate unique keys
    check(len(set(l)), len(l), "Keys are not unique")

    # validate parameters for each location
    for lsub in l:
        d = data[lsub]

        # validate parameters are present
        keyCheck = [x in d.keys() for x in lKeys]
        missing = list(compress(lKeys, [not x for x in keyCheck]))
        check(
            True,
            all(keyCheck),
            f"Key(s) ({', '.join(missing)}) are missing for location '{lsub}'",
        )

        # validate lon + lat are strings that can be numbers
        [checkStrFloat(d[k]) for k in ["lon", "lat"]]

        # validate units are valid
        check(
            True, d["units"] in units, f"'{d['units']}' is not a valid units parameter"
        )


def verifySchedule(data):
    l = list(data.keys())
    v = list(data.values())
    days = [
        "Monday",
        "Tuesday",
        "Wednesday",
        "Thursday",
        "Friday",
        "Saturday",
        "Sunday",
    ]

    # validate unique keys
    check(len(set(l)), len(l), "Keys are not unique")

    # validate existing keys
    keyCheck = [x in days for x in l]
    invalid = list(compress(l, [not x for x in keyCheck]))
    check(True, all(keyCheck), f"{', '.join(invalid)} are not valid keys")

    # validate existing values
    [datetime.strptime(x, "%I:%M%p") for x in v]


if __name__ == "__main__":
    verifyLocations(getData("locations.json"))
    verifySchedule(getData("schedule.json"))
