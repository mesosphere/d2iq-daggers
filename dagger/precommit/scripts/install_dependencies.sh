#!/usr/bin/env bash

set -eox pipefail

# Install dependencies

apt-get update

apt-get install -y --no-install-recommends bash git python3 python3-pip python3-dev build-essential

pip3 install --upgrade pre-commit gitlint

rm -rf /var/lib/apt/lists/*
