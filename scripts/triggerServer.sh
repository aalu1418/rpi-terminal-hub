#!/bin/bash

# exit when any command fails
set -e

# read server address from .env
SERVER=$(grep SERVER= <.env | cut -d '=' -f2)

# default server if not specified
if [ "$SERVER" = "" ]; then
  SERVER="localhost"
fi

BRANCH=$(git rev-parse --abbrev-ref HEAD)

# only trigger if push occurs on the 'main' branch
if [ "$BRANCH" = "main" ]; then
  curl -X POST "$SERVER:5000/pre-push"
fi
