#!/bin/bash

# exit when any command fails
set -e

# read server address from .env
SERVER=$(grep SERVER= <.env | cut -d '=' -f2)

# default server if not specified
if [ "$SERVER" = "" ]; then
  SERVER="localhost"
fi

sleep 5 && curl -X POST "$SERVER:5000/pull" > /dev/null &


# sleep 10 && curl -X POST "$SERVER:5000/reboot" &
