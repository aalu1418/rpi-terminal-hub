import requests, re

class TTC:
    def __init__(self):
        self.url = "https://www.ttc.ca/Service_Advisories/all_service_alerts.jsp"

    def fetch(self):
        res = requests.get(self.url)
        self.data = res.text
        self.parse()

    def parse(self):
        # parse webpage + filter out non relevant alerts
        self.data = re.findall(r'<div class="alert-content"><p class="veh-replace">.*</p></div></div>', self.data)[0]
        self.data = re.split(r'<div class="alert-content"><p class="veh-replace">', self.data)[1:]
        self.data = [re.sub(r':.*', '', v) for v in self.data if "Regular service has resumed" not in v and "Elevator" not in v]

        # filter  for streetcars + subways
        temp = [v[:6] for v in self.data if "Line" in v] # only pulls `Line #`
        temp = temp + [' '.join(v.split(" ")[0:2]) for v in self.data if re.search(r'5[0-9]{2}[a-zA-Z]{0,1}.*', v)]
        self.data = temp

if __name__ == '__main__':
    ttc = TTC()
    ttc.fetch()
    print(ttc.data)
