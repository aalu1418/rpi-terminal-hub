#!/bin/bash

# exit when any command fails
set -e

# install packages
sudo apt-get update
sudo apt-get install -y git python3-pip pigpio vim

# set timezone to US/Eastern time
sudo timedatectl set-timezone America/New_York

# install python packages
pip3 install pigpio ircodec flask requests python-dotenv feedparser

# pull repo
git clone https://github.com/aalu1418/rpi-terminal-hub

# setup cron job
crontab -l | { cat; echo "@reboot sleep 60 && /usr/bin/python3 /home/pi/rpi-terminal-hub/main.py >> ~/cron.log 2>&1"; } | crontab -
