import json, os

# class for generating bash commands
class Cmd:
    def __init__(self, s=0):
        self.change(s)

    def change(self, s):
        if not isinstance(s, int):
            raise TypeError("Not an integer")
        self.s = s
        self.update()

    def update(self):
        self.pull = f"sleep {self.s} && cd /home/pi/rpi-terminal-hub && git pull"
        self.rmLog = f"sleep {self.s+5} && rm /home/pi/cron.log"
        self.reboot = f"sleep {self.s+10} && sudo reboot"


# fetching json data from data folder
def getData(name):
    dirname = os.path.dirname(__file__)  # use filepath to utils file
    datapath = os.path.realpath(os.path.join(dirname, "../data"))
    f = open(f"{datapath}/{name}")
    schedule = json.load(f)
    f.close()

    return schedule


if __name__ == "__main__":
    c = Cmd(s=5)
    print(c.__dict__)
    c.change(0)
    print(c.__dict__)
