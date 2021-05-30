#!/bin/bash

# exit when any command fails
set -e

# read server address from .env
SERVER=$(grep SERVER= <.env | cut -d '=' -f2)

# default server if not specified
if [ "$SERVER" = "" ]; then
  SERVER="localhost"
fi

curl -X POST "$SERVER:5000/pull"

sleep 5 && curl -X POST "$SERVER:5000/reboot"
