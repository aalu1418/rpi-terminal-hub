#!/bin/bash

# install packages
sudo apt-get update
sudo apt-get install -y git python3-pip pigpio vim

# install python packages
pip3 install pigpio ircodec flask requests python-dotenv

# pull repo
git pull https://github.com/aalu1418/rpi-terminal-hub

# setup cron job
crontab -l | { cat; echo "@reboot sleep 60 && /usr/bin/python3 /home/pi/rpi-terminal-hub/main.py --server >> ~/cron.log 2>&1"; } | crontab -
